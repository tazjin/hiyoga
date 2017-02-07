package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/polydawn/meep"
	"github.com/tazjin/hiyoga/util"
	"time"
)

const LOGIN_URL string = "https://www.hiyoga.no/webapi/auth/login"
const REQUEST_CONTENT_TYPE = "application/json;charset=UTF-8"
const COOKIE_NAME = "Auth-SatsElixia"

// Login request type as expected by HiYoga API
type loginRequest struct {
	UserName string
	Password string
}

type loginResponse struct {
	LoginSuccess  bool
	StatusMessage string
	UserId        string `json:"userId"`
}

type LoginError struct {
	meep.TraitTraceable
	meep.TraitAutodescribing
	meep.TraitCausable

	Response string
}

func AuthenticatedGet(url string) (resp *http.Response, err error) {
	token := getCredentials()

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.AddCookie(&http.Cookie{
		Name:  COOKIE_NAME,
		Value: token.Token,
	})

	return client.Do(req)
}

// Attempts to load existing credentials from disk and fetches new ones if necessary
func getCredentials() *util.HiyogaToken {
	config, err := util.LoadConfig()

	if err != nil {
		util.Fail(err)
	}

	if config.Token != nil && config.Token.TokenStillValid() {
		return config.Token
	}

	// If we get here we need a new token!
	if config.Username == nil || config.Password == nil {
		util.Fail(fmt.Errorf("Username and password not configured"))
	}

	token, err := performLoginCall(&loginRequest{
		UserName: *config.Username,
		Password: *config.Password,
	})

	if err != nil {
		util.Fail(err)
	}

	// And update it in the configuration
	config.Token = token
	util.WriteConfig(config)

	return token
}

func performLoginCall(request *loginRequest) (*util.HiyogaToken, error) {
	var l loginResponse
	postBody, _ := json.Marshal(*request)
	resp, err := http.Post(LOGIN_URL, REQUEST_CONTENT_TYPE, bytes.NewBuffer(postBody))

	if err != nil {
		util.Fail(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 || !l.LoginSuccess {
		util.Fail(fmt.Errorf("Login failed: %s (%d)\n", string(body), resp.StatusCode))
	}

	err = json.Unmarshal(body, &l)

	if err != nil {
		util.Fail(err)
	}

	token := findAuthCookie(resp)

	if token == "" {
		util.Fail(fmt.Errorf("No authentication token found in login response"))
	}

	credentials := &util.HiyogaToken{
		Timestamp: time.Now(),
		Token:     token,
	}

	return credentials, nil
}

func findAuthCookie(resp *http.Response) string {
	for _, c := range resp.Cookies() {
		if c.Name == COOKIE_NAME {
			return c.Value
		}
	}

	return ""
}
