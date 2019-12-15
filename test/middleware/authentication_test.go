package middleware

import (
	"encoding/json"
	m "github.com/saeedafshari8/authenticator/middleware"
	"github.com/saeedafshari8/authenticator/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func getLoginPOSTPayload() string {
	login := m.Login{
		Username: "admin",
		Password: "admin",
	}
	result, err := json.Marshal(login)
	if err != nil {
		panic(err)
	}
	return string(result)
}

func TestUnauthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := test.GetRouter()

	loginPOSTPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/v1/echo", strings.NewReader(loginPOSTPayload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPOSTPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}

	p, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fail()
	}

	var response m.HttpStatusResponse
	err = json.Unmarshal(p, &response)
	if err != nil || response.Message != "cookie token is empty" {
		t.Fail()
	}
}

func TestOpenAPI(t *testing.T) {
	w := httptest.NewRecorder()

	r := test.GetRouter()

	loginPOSTPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/v1/login", strings.NewReader(loginPOSTPayload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPOSTPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestLoginAndGetToken(t *testing.T) {
	w := httptest.NewRecorder()

	r := test.GetRouter()

	loginPOSTPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/v1/login", strings.NewReader(loginPOSTPayload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPOSTPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	p, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fail()
	}

	var response m.LoginResponse
	err = json.Unmarshal(p, &response)
	if err != nil || response.Token == "" {
		t.Fail()
	}
}
