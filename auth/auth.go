package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hiyoga/util"
	"io/ioutil"
	"net/http"

	"github.com/polydawn/meep"
	"os"
)

const LOGIN_URL string = "https://www.hiyoga.no/webapi/auth/login"
const REQUEST_CONTENT_TYPE = "application/json;charset=UTF-8"
const COOKIE_NAME = "Auth-SatsElixia"

type LoginRequest struct {
	UserName string
	Password string
}

type loginResponse struct {
	LoginSuccess  bool
	StatusMessage string
	UserId        string `json:"userId"`
}

type Credentials struct {
	UserId string
	Token  string
}

type LoginError struct {
	meep.TraitTraceable
	meep.TraitAutodescribing
	meep.TraitCausable

	Response string
}

func LoginAndStoreCredentials(request *LoginRequest) {
	creds, err := performLoginCall(request)

	if err != nil {
		util.Fail(err)
	}

	json, _ := json.Marshal(creds)

	ioutil.WriteFile(getLoginFileLocation(), json, 0644)
	fmt.Printf("Stored credentials for user %s\n", creds.UserId)
}

func performLoginCall(request *LoginRequest) (*Credentials, error) {
	var l loginResponse
	postBody, _ := json.Marshal(*request)
	resp, err := http.Post(LOGIN_URL, REQUEST_CONTENT_TYPE, bytes.NewBuffer(postBody))

	if err != nil {
		util.Fail(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &l)

	if resp.StatusCode != 200 || !l.LoginSuccess {
		util.Fail(fmt.Errorf("Login failed: %s (%d)\n", string(body), resp.StatusCode))
	}

	if err != nil {
		util.Fail(err)
	}

	token := findAuthCookie(resp)

	if token == "" {
		util.Fail(fmt.Errorf("No authentication token found in login response"))
	}

	credentials := &Credentials{
		UserId: l.UserId,
		Token:  token,
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

func getLoginFileLocation() string {
	home := os.Getenv("HOME")
	return home + "/.hiyoga"
}
