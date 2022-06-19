package main

import (
	"net/http"

	"github.com/colt005/bulkIt"
)

func main() {
	c := &http.Client{}
	maxThreads := 10
	b := bulkIt.NewBulkIt(c, maxThreads)

	b.SaveFilesByUrls([]string{"https://www.sample-videos.com/img/Sample-jpg-image-1mb.jpg"}, "./")
}
