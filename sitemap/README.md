# Sitemap builder

## What does it do?
Starting from the root directory creates a sitemap of links.  
Saves the map in an xml file. Check [Description](Description.md) for more details.

## How to use?
1. Download and cd into sitemap folder
2. `go run .` Options are `go run . -h`:
    * site: Full path to site e.g. `http://www.example.com`
    * depth: Depth from root node. (default : 3)
    * count: Max number of urls to explore. (default : 32)
3. Results are saved in `sitemap.xml`

## Tasks
- [x] Modularize link-parser
- [x] Fetch web-pages and get all links
- [x] Use BFS to list out all links in order of visit