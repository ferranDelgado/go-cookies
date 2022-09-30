package pkg

import (
	json "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go StartApp()
	waitServer()
	run := m.Run()
	println("Test end")
	os.Exit(run)
}

func Test_WithoutCookieItShouldReturnNotAuthorize(t *testing.T) {
	code, _ := getCall("/v1/auth")
	if code != http.StatusUnauthorized {
		t.Errorf("Incorrect Status code %v", code)
	}
}

func Test_AfterLoginTheresAValidCookie(t *testing.T) {
	code, _ := postCall("/login", &gin.H{
		"user":     "root_user",
		"password": "root_password",
	})
	assert.Equal(t, code, http.StatusOK)
	code, _ = getCall("/v1/auth")
	assert.Equal(t, code, http.StatusOK)
}

func waitServer() {
	code, _ := getCall("/health")
	for code != http.StatusOK {
		code, _ = getCall("/health")
	}
}

func getCall(path string) (int, string) {
	return call("GET", path, nil)
}

func postCall(path string, data *gin.H) (int, string) {
	return call("POST", path, data)
}

var jar, _ = cookiejar.New(nil)

func call(method string, path string, data *gin.H) (int, string) {
	fullUrl := fmt.Sprintf("http://localhost:%d%v", 8080, path)
	var httpClient = http.Client{
		Timeout: time.Duration(24) * time.Hour,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}
	jsonData, err := json.Marshal(data)

	req, _ := http.NewRequest(method, fullUrl, strings.NewReader(string(jsonData)))

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Print(err.Error())
		return 0, err.Error()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	u := &url.URL{
		Scheme: "http",
		Host:   "example.com",
	}
	jar.SetCookies(u, resp.Cookies())
	return resp.StatusCode, string(body)
}

type tokenResponse struct {
	Token string `json:"token"`
}

func unmarshal[K any](value string) K {
	var result = new(K)
	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return *result
}
