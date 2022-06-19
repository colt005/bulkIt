# bulkIt [![GoDoc](https://godoc.org/github.com/natefinch/godocgo?status.png)]([https://godoc.org/github.com/natefinch/godocgo](https://pkg.go.dev/github.com/colt005/bulkIt@v1.0.0))


## About The Package

A Go package that allows you to download multiple files concurrently.

_Get Files, One Goroutine at a time_

### Import package
```sh
go get github.com/colt005/bulkIt
```

```go
import "github.com/colt005/bulkIt"
```

## Usage

```go
    c := &http.Client{}
    maxThreads := 50
    b := bulkIt.NewBulkIt(c,maxThreads)
```

## Example

```go
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

```

## Functions

### ```func func bulkIt.NewBulkIt(c *http.Client, maxThreads int) *bulkIt.bulkItClient```
Create a new bulkItClient.

- Params
    - c - The bulkIt package's http client for sending out http requests.If you specify nil, the package will utilise the default http client. 
    - maxThreads - The maximum number of go routines that may be executed simultaneously to send out http requests. If the value is sepcified as 0, It will default to 30. 

### ```func (*bulkItClient) SaveFilesByUrls(urls []string,path string)```

SaveFilesByUrls takes slice of urls and the path where the files need to be saved.


