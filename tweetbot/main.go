package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	creds map[string]string
)

func init() {
	data, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &creds); err != nil {
		panic(err)
	}
}

func main() {
	accessToken, err := getAccessToken()
	check(err)

	fmt.Println(accessToken)
}

func getAccessToken() (string, error) {
	client := &http.Client{}
	reqBody := strings.NewReader("grant_type=client_credentials")
	req, _ := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", reqBody)

	var bearerToken = getBearerToken()
	req.Header.Set("Authorization", "Basic "+bearerToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	check(err)
	var payload map[string]string
	err = json.Unmarshal(data, &payload)
	check(err)

	if payload["token_type"] != "bearer" {
		return "", errors.New("Invalid Auth Token")
	}
	return payload["access_token"], nil
}

func getBearerToken() string {
	consumerKey := url.QueryEscape(creds["API_KEY"])
	consumerSecret := url.QueryEscape(creds["API_SECRET"])
	bearer := consumerKey + ":" + consumerSecret
	return base64.StdEncoding.EncodeToString([]byte(bearer))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
