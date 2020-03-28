package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/wangwalton/gocrawler/contracts"
	"github.com/wangwalton/gocrawler/scraper/html"
	"github.com/wangwalton/gocrawler/scraper/queue"
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

func ProcessURL(client contracts.URLQueueClient, scrape_url string) {
	body, _ := Get(scrape_url)
	url_object, _ := url.Parse(scrape_url)
	links := html.Extractor(body, url_object)
	fmt.Printf("Processed %s found %d links\n", scrape_url, len(links))
	for l := range links {
		sj := contracts.ScraperJob{Url: l}
		queue.Enqueue(client, &sj)
	}
}
