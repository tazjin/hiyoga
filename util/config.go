package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

// A HiYoga JWT and the timestamp at which it was retrieved
type HiyogaToken struct {
	Token     string    `json:"token"`
	Timestamp time.Time `json:"timestamp"`
}

// Configuration for HiYoga CLI is used to determine how logging in should be managed.
type HiyogaConfig struct {
	Username *string      `json:"username"`
	Password *string      `json:"password"`
	Token    *HiyogaToken `json:"token"`
}

// load the configuration from the user's home directory
func LoadConfig() (HiyogaConfig, error) {
	var c HiyogaConfig
	content, err := ioutil.ReadFile(getConfigFileLocation())

	if err != nil {
		return c, err
	}

	err = json.Unmarshal(content, &c)

	if err != nil {
		return c, err
	}

	return c, nil
}

// write updated configuration to the user's home directory
func WriteConfig(config HiyogaConfig) error {
	content, _ := json.Marshal(config)

	return ioutil.WriteFile(getConfigFileLocation(), content, 0600)
}

func getConfigFileLocation() string {
	return os.Getenv("HOME") + "/.hiyoga"
}

// Check whether a JWT should still be considered valid (less than 60m old)
func (t *HiyogaToken) TokenStillValid() bool {
	return t.Timestamp.Add(60 * time.Minute).After(time.Now())
}
