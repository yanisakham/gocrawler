package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/wangwalton/gocrawler/queue"
	"github.com/wangwalton/gocrawler/scraper/html"
)

func Get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func ProcessURL(q queue.Queue, scrape_url string) {
	body, _ := Get(scrape_url)
	url_object, _ := url.Parse(scrape_url)
	links := html.Extractor(body, url_object)
	fmt.Printf("Processed %s found %d links\n", scrape_url, len(links))
	for l := range links {
		q.Enqueue(l)
	}
}
