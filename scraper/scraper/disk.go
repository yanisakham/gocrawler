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

	escaped_path := strings.ReplaceAll(u.Path, "/", "_")
	filename := folder + "/" + escaped_path + ".html"
	zap.S().Debugf("Writing file %", filename)
	err := ioutil.WriteFile(filename, document, 0700)
	if err != nil {
		zap.S().Warnf("writing error %v", err)
	}
}
