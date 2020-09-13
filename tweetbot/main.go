package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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
	rand.Seed(time.Now().UnixNano())
}

// TwitterUser is the type stored in the localstore
type TwitterUser struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

// Represent Twitter user in a single line
func (twu TwitterUser) String() string {
	return fmt.Sprintf("%v @ %v - %v", twu.Name, twu.ScreenName, twu.UserID)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Incorrect usage. Specify fetch or choose")
		fmt.Println("tweetbot fetch <tweet-id>")
		fmt.Println("tweetbot choose <tweet-id> <no-of-users-to-pick>")
		return
	}

	switch os.Args[1] {
	case "fetch":
		if len(os.Args) < 3 {
			fmt.Println("Need tweet-id")
			return
		}
		fetch(os.Args[2])
	case "choose":
		if len(os.Args) < 3 {
			fmt.Println("Need tweet-id")
			return
		}
		var picks = 1
		if len(os.Args) < 4 {
			fmt.Println("Did not specify how many people to pick. Default - 1")
		} else {
			var err error
			picks, err = strconv.Atoi(os.Args[3])
			check(err)
		}
		choose(os.Args[2], picks)
	default:
		fmt.Printf("Invalid choice %v, Only fetch or choose", os.Args[1])
	}
}

func fetch(tweetID string) {
	// Fetch once
	accessToken, err := getAccessToken()
	check(err)
	retweeters := getRetweeters(tweetID, accessToken)
	var fileName = fmt.Sprintf("tweeters-%s.json", tweetID)
	saveHistory(retweeters, fileName)

	// List them all by usernames
	fmt.Println("Retweeters -")
	tweeters := loadHistory(fileName)
	var mlt, n = 100, len(tweeters)

	for i := 0; i*mlt < n; i++ {
		var r = mlt * (i + 1)
		if r > n {
			r = n
		}
		chunk := tweeters[mlt*i : r]
		users := getUsers(chunk, accessToken)
		for j := range users {
			fmt.Printf("%s @ %s\n", users[j].Name, users[j].ScreenName)
		}
	}
}

func choose(tweetID string, picks int) {
	// Fetch once
	accessToken, err := getAccessToken()
	check(err)
	retweeters := getRetweeters(tweetID, accessToken)
	var fileName = fmt.Sprintf("tweeters-%s.json", tweetID)
	saveHistory(retweeters, fileName)

	tweeters := loadHistory(fileName)
	n := len(tweeters)
	if n <= 0 {
		return
	}
	if picks > n {
		picks = n
	}
	indexes := []int{}
	var set = map[int]bool{}
	for picks > 0 {
		ix := rand.Intn(n)
		if _, ok := set[ix]; ok {
			continue
		}
		indexes = append(indexes, ix)
		picks--
	}

	var winners []string
	for _, e := range indexes {
		winners = append(winners, tweeters[e])
	}

	// Add paging support
	winUsers := getUsers(winners, accessToken)
	fmt.Println("Winners are -")
	for i := range winUsers {
		fmt.Printf("%s @ %s\n", winUsers[i].Name, winUsers[i].ScreenName)
	}
}

// Loads all retweeters ids stored in file
func loadHistory(fileName string) []string {
	data, err := ioutil.ReadFile(fileName)
	check(err)
	var history []string
	json.Unmarshal(data, &history)
	return history
}

// Save new retweeter ids to file along with old ones
func saveHistory(retweeters []string, fileName string) {

	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		f, err := os.Create(fileName)
		check(err)
		f.Close()
	}

	data, err := ioutil.ReadFile(fileName)
	check(err)
	var history []string
	json.Unmarshal(data, &history)

	// filter out repeated entrires
	// done by json
	seen := map[string]bool{}
	for i := range history {
		seen[history[i]] = true
	}

	for i := range retweeters {
		if _, ok := seen[retweeters[i]]; !ok {
			history = append(history, retweeters[i])
			seen[retweeters[i]] = true
		}
	}

	f, err := os.Create(fileName)
	check(err)
	rtB, _ := json.Marshal(history)
	f.Write(rtB)
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

	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweeters/ids.json")
	req, _ := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("id", id)
	q.Add("trim_user", "true")
	q.Add("count", "100")
	q.Add("stringify_ids", "true")
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

	var users []string

	// for k, v := range payload {
	// 	fmt.Println(k, v)
	// }

	for i := range payload["ids"].([]interface{}) {
		user := payload["ids"].([]interface{})[i]
		users = append(users, user.(string))
	}

	// Debugging
	// ck := payload[0]
	// for k := range ck {
	// 	fmt.Println(k)
	// }
	// fmt.Println(ck["user"])

	// for i := range payload {
	// 	user := payload[i]
	// 	users = append(users, user["user"].(map[string]interface{})["id_str"].(string))
	// }

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
