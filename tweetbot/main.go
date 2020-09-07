package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	fmt.Println(creds)
}
