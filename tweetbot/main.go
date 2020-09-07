package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	resp, err := http.Get("https://api.twitter.com/1.1/statuses/show.json?id=210462857140252672")
	check(err)
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	fmt.Println(string(body))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
