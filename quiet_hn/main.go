package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"soul-sapphire/quiet_hn/hn"
)

type concItem struct {
	article hn.Item
	loc     int
}

type storyItem struct {
	story item
	loc   int
}

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	fmt.Println("Serving on port 3000")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	var client hn.Client
	// client.ItemCache = make(map[int]hn.Item)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ids, err := client.TopItems()
		fetched := time.Now().Sub(start)
		fmt.Println("Fetch time is:", fetched)

		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		// fmt.Println(client.ItemCache)
		var stories = map[int]item{}
		itemChan := make(chan concItem)
		parseStart := time.Now()
		for _, id := range ids {
			go func(itemId int) {
				hnItem, err := client.GetItem(itemId)
				if err != nil {
					return
				}
				itemChan <- concItem{hnItem, itemId}
			}(id)
		}

		for i := 0; i < numStories; {
			fullItem := <-itemChan
			hnItem := fullItem.article
			article := parseHNItem(hnItem)
			if isStoryLink(article) {
				stories[fullItem.loc] = article
				i++
			}
		}

		// Put away all extra items into a blackhole
		go func() {
			for {
				_ = <-itemChan
			}
		}()

		var sortStory []storyItem

		for k, v := range stories {
			sortStory = append(sortStory, storyItem{
				story: v,
				loc:   k,
			})
		}

		sort.Slice(sortStory, func(i, j int) bool {
			return sortStory[i].loc < sortStory[j].loc
		})

		var finalList []item
		for i := range sortStory {
			finalList = append(finalList, sortStory[i].story)
		}

		parseTime := time.Now().Sub(parseStart)
		fmt.Println("Parse time:", parseTime)
		data := templateData{
			Stories: finalList,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
