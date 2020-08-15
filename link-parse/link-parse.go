package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

type parsedTag struct {
	Link string `json:"Href"`
	Data string `json:"Text"`
}

func main() {
	// Get html name through CLA
	// Allow functionality to get file from the web
	var fileName = flag.String("html", "ex3.html", "The file to parse")
	var urlName = flag.String("url", "", "The url to get the file from")
	flag.Parse()

	var fileLoc = *fileName
	var overHTTP = false
	if *urlName != "" {
		resp, err := http.Get(*urlName)
		if err != nil {
			fmt.Println("Cannot fetch file from that url")
			return
		}
		defer resp.Body.Close()
		f, _ := os.Create("tmp.html")
		io.Copy(f, resp.Body)
		f.Close()

		fileLoc = "tmp.html"
		overHTTP = true
	}

	reader, err := os.Open(fileLoc)
	if err != nil {
		fmt.Println("Could not read file")
		return
	}

	// Got the io Reader, Get the parse root
	root, perr := html.Parse(reader)
	if perr != nil {
		fmt.Println("Error while parsing html")
		return
	}
	reader.Close()

	// DFS parse the file from root to get all a tags
	var tags = []*html.Node{}

	tags = dfs(root, tags)

	var result = make([]parsedTag, len(tags))

	// Get all data within the a tags
	for i, tag := range tags {
		link := tag.Attr[0].Val
		tagData := extract(tag)
		tmp := parsedTag{
			Link: link,
			Data: tagData,
		}
		result[i] = tmp
	}

	// Convert lists to JSON
	serial, jerr := json.MarshalIndent(result, "", "  ")
	if jerr != nil {
		fmt.Println("Could not convert to JSON")
		return
	}

	outFile, err := os.Create("parsed.json")
	if err != nil {
		fmt.Println("Cannot open file to write results")
		return
	}
	defer outFile.Close()

	_, werr := outFile.Write(serial)
	if werr != nil {
		fmt.Println("Could not write to file")
		return
	}

	outFile.Sync()

	if overHTTP {
		os.Remove(fileLoc)
	}
}

// dfs recursively gets all the anchor tags within the root
func dfs(root *html.Node, tags []*html.Node) []*html.Node {
	if root == nil {
		return nil
	}

	if root.Type == html.ElementNode && root.Data == "a" {
		tags = append(tags, root)
	}

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		tags = dfs(child, tags)
	}

	return tags
}

// extract gets all the text within a tag, recursively
func extract(root *html.Node) string {
	if root == nil {
		return ""
	}

	var tmp string
	for child := root.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode && child.Type != html.CommentNode {
			// fmt.Println(child, child.Data)
			tmp += child.Data
		}
		tmp += extract(child)
	}

	return tmp
}
