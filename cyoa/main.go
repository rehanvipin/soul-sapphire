package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type storyArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []option `json:"options"`
}

func main() {
	data, err := ioutil.ReadFile("gopher.json")
	if err != nil {
		fmt.Println("Cannot read file")
		return
	}
	var flow = map[string]storyArc{}
	jsonerr := json.Unmarshal(data, &flow)
	if jsonerr != nil {
		fmt.Println("Could not parse JSON file")
		return
	}

	// Serve pages through http
	http.HandleFunc("/", homePage)

	// Special for intro
	stories := storyReader(flow, "story.html")

	fmt.Println("Serving on port 8040")
	http.ListenAndServe("localhost:8040", stories)
}

func homePage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s\n", "<h1>Welcome to cyoa</h1>")
}

func storyReader(available map[string]storyArc, temple string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := strings.Split(req.URL.Path, "/")
		if path[1] == "" || path[1] == "/" {
			path[1] = "intro"
		}
		if arc, ok := available[path[1]]; ok {
			tmpl := template.Must(template.ParseFiles(temple))
			tmpl.Execute(w, arc)
		} else {
			http.Error(w, "404 - Could not find that story.", http.StatusNotFound)
			log.Printf("%s %v\n", "Story not found", req.URL.Path)
		}
	}
}
