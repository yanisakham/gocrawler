package scraper

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/wangwalton/gocrawler/contracts"
	"go.uber.org/zap"
)

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func ProcessURL(client contracts.URLQueueClient, scrape_url string) {
	body, _ := Get(scrape_url)
	url_object, _ := url.Parse(scrape_url)
	links := ExtractorHTML(string(body), *url_object)
	zap.S().Debugf("Processed %s found %d links", scrape_url, len(links))

	go WriteFile(body, url_object)
	for l := range links {
		sj := contracts.ScraperJob{Url: l}
		go Enqueue(client, &sj)
	}
}
