package scraper

import (
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ExtractorHTML(h string, u url.URL) map[string]bool {
	doc, err := html.Parse(strings.NewReader(h))
	if err != nil {
		log.Fatal(err)
	}

	links := make(map[string]bool)
	extractorHelper(doc, u, links)
	return links
}

const https = "https:"

func extractorHelper(n *html.Node, u url.URL, links map[string]bool) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {

				var parsed_url string
				if strings.HasPrefix(a.Val, "//") {
					parsed_url = https + a.Val
				} else if strings.HasPrefix(a.Val, "/") {
					u.Path = a.Val
					parsed_url = u.String()
				} else if a.Val == "#" {
					continue
				} else {
					parsed_url = a.Val
				}
				links[parsed_url] = true
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractorHelper(c, u, links)
	}
}
