package main

import (
	"errors"
	"net/http"

	"golang.org/x/net/html"
)

// LinkFetch gets all the links from a webpage given the URL
func LinkFetch(URL string) ([]string, error) {

	resp, err := http.Get(URL)
	if err != nil {
		return nil, errors.New("Cannot fetch data from that url")
	}
	defer resp.Body.Close()

	// Got the io Reader, Get the parse root
	root, perr := html.Parse(resp.Body)
	if perr != nil {
		return nil, errors.New("Cannot parse html")
	}

	// DFS parse the file from root to get all a tags
	var tags = []*html.Node{}

	tags = dfs(root, tags)

	var result = make([]string, len(tags))

	// Get all data within the a tags
	for i, tag := range tags {
		link := tag.Attr[0].Val
		result[i] = link
	}

	return result, nil
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
