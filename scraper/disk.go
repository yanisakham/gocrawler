package scraper

import (
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
)

func WriteFile(document []byte, u *url.URL) {

	folder := "data/" + u.Hostname()
	os.MkdirAll(folder, 0700)

	var path string
	if len(u.Path) > 0 {
		path = u.Path[1:]
	} else {
		path = u.Path
	}
	escaped_path := strings.ReplaceAll(path, "/", "_")
	filename := folder + "/" + escaped_path + ".html"
	zap.S().Debugf("Writing file %s", filename)
	err := ioutil.WriteFile(filename, document, 0700)
	if err != nil {
		zap.S().Warnf("writing error %v", err)
	}
}
