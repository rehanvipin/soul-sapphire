package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Link is the intermediary struct for XML
type Link struct {
	XMLName xml.Name `xml:"url"`
	URL     string   `xml:"loc"`
}

// Result is the format required for XML
type Result struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Links   []*Link  `xml:url>loc`
}

func main() {
	// Set URL through CLAs
	// Also set maximum depth and count
	siteFlag := flag.String("site", "http://www.example.com", "The url to build the sitemap for")
	dipFlag := flag.Int("depth", 3, "Maximum depth of search")
	countFlag := flag.Int("count", 32, "Max no of sub-urls to explore")
	flag.Parse()

	siteRoot := *siteFlag
	maxDepth := *dipFlag
	maxCount := *countFlag

	fmt.Println("Creating sitemap for", siteRoot)
	siteMap := bfs(siteRoot, maxDepth, maxCount)

	// fmt.Println(siteMap)

	collection := []*Link{}
	for i := range siteMap {
		collection = append(collection, &Link{
			URL: siteMap[i],
		})
	}

	formatted := &Result{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Links: collection,
	}
	out, xerr := xml.MarshalIndent(formatted, " ", "  ")
	check(xerr)

	f, err := os.Create("sitemap.xml")
	check(err)
	f.Write([]byte(xml.Header))
	f.Write(out)
	f.Close()
}

func bfs(URL string, maxDepth, maxCount int) []string {
	u, err := url.Parse(URL)
	check(err)
	scheme := u.Scheme
	hostOG := u.Host

	var visited = make(map[string]bool)
	var start = URL
	var bfsq = []string{start}
	var depths = []int{0}
	var history = []string{}

	var count = 0

	for len(bfsq) != 0 {
		start = bfsq[0]
		depth := depths[0]
		depths = depths[1:]

		bfsq = bfsq[1:]
		visited[start] = true

		if depth > maxDepth {
			continue
		} else if count > maxCount {
			break
		}

		links, gerr := LinkFetch(start)
		if gerr != nil {
			fmt.Printf("Skipping %s\n", start)
			continue
		}
		fmt.Printf("Checking out %s at level %d\n", start, depth)
		count++

		for _, link := range links {
			if strings.HasPrefix(link, "/") {
				link = scheme + "://" + hostOG + link
			}

			purl, perr := url.Parse(link)
			if perr != nil || purl.Host != hostOG {
				continue
			}

			if _, ok := visited[link]; !ok {
				bfsq = append(bfsq, link)
				depths = append(depths, depth+1)
			}
		}
		history = append(history, start)
	}

	return history
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
