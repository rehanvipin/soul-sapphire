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

// Populate the API credentials key-value pair
func init() {
	data, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &creds); err != nil {
		panic(err)
	}
}

// TwitterUser is the type stored in the localstore
type TwitterUser struct {
	UserID, Name, ScreenName string
}

// Represent Twitter user in a single line
func (twu TwitterUser) String() string {
	return fmt.Sprintf("%v @ %v - %v", twu.Name, twu.ScreenName, twu.UserID)
}

func main() {
	accessToken, err := getAccessToken()
	check(err)

	var tweetID = "1303359966741508096"
	retweeters := getRetweeters(tweetID, accessToken)
	RTUsers := getUsers(retweeters, accessToken)

	for i := range RTUsers {
		fmt.Println(RTUsers[i])
	}
}

// Use user-ids to get TwitterUser objects
func getUsers(userIDs []string, accessToken string) []TwitterUser {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.twitter.com/1.1/users/lookup.json")
	req, _ := http.NewRequest("GET", url, nil)

	// fmt.Println(userIDs)

	userIDlist := ""
	for i := range userIDs {
		if i > 0 {
			userIDlist += ","
		}
		userIDlist += userIDs[i]
	}

	// fmt.Println(userIDlist)

	q := req.URL.Query()
	q.Add("user_id", userIDlist)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	check(err)

	// fmt.Println(string(data))

	var payload []map[string]interface{}
	err = json.Unmarshal(data, &payload)
	check(err)

	var users []TwitterUser

	for i := range payload {
		user := payload[i]
		// fmt.Println(user["id_str"], user["name"], user["screen_name"])
		users = append(users, TwitterUser{
			UserID:     user["id_str"].(string),
			Name:       user["name"].(string),
			ScreenName: user["screen_name"].(string),
		})
	}

	return users
}

// Return array of user-ids of all retweeters
func getRetweeters(id, accessToken string) []string {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%s.json", id)
	req, _ := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("trim_user", "true")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	check(err)

	var payload []map[string]interface{}
	err = json.Unmarshal(data, &payload)
	check(err)

	var users []string

	// Debugging
	// ck := payload[0]
	// for k := range ck {
	// 	fmt.Println(k)
	// }
	// fmt.Println(ck["user"])

	for i := range payload {
		user := payload[i]
		users = append(users, user["user"].(map[string]interface{})["id_str"].(string))
	}

	return users
}

// Return all the properties of a tweet with a given id
func getTweet(id, accessToken string) map[string]interface{} {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitter.com/1.1/statuses/show.json", nil)

	q := req.URL.Query()
	q.Add("id", id)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	check(err)
	var payload map[string]interface{}
	err = json.Unmarshal(data, &payload)
	check(err)

	return payload
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
