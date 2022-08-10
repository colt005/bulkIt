package bulkIt

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
	"github.com/zenthangplus/goccm"
)

type bulkItClient struct {
	client     *httpClient
	maxThreads int
}

type httpClient struct {
	client *http.Client
}

/// Create an FileSave object which takes in an http Client to make requests and an integer value which controls the number of threads to be fired.
// If `c` is nil, default http.Client will be assigned.
// If `maxThreads` is 0 then `maxThreads` will default to 30
func NewBulkIt(c *http.Client, maxThreads int) *bulkItClient {
	if c == nil {
		c = &http.Client{}
	}

	if maxThreads == 0 {
		maxThreads = 30
	}

	h := &httpClient{
		client: c,
	}

	return &bulkItClient{
		client:     h,
		maxThreads: maxThreads,
	}
}

func (c *httpClient) GET(url string) (*http.Response, error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (im *bulkItClient) SaveFilesByUrls(urls []string, path string) {
	startTime := time.Now()

	c := goccm.New(im.maxThreads)
	counter := 0
	for i := 0; i < len(urls); i++ {
		c.Wait()
		go func(i int64) {
			err := worker(*im.client, urls[i], path)
			if err != nil {
				logger.Error("Failed to download url : ", urls[i])
			}
			if err == nil {
				counter += 1
			}
			fmt.Printf("Downloaded %d of %d \n", i+1, len(urls))
			c.Done()
		}(int64(i))
	}
	c.WaitAllDone()
	fmt.Printf("%d of %d files downloaded successfully", counter, len(urls))
	logger.Info("Time Taken : " + time.Since(startTime).String())
}

func worker(client httpClient, url string, fPath string) (err error) {
	if url == "" {
		return
	}
	// fmt.Println("Downloading..." + url)
	resp, err := client.GET(url)
	if err != nil {
		fmt.Println(err.Error())
		if resp != nil {
			resp.Body.Close()
		}
		return errors.New(`failed to download. Please check url`)
	}
	defer resp.Body.Close()
	// fmt.Println(resp.Header)

	if resp.StatusCode != http.StatusOK {
		return errors.New(`failed to download. Please check url`)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return errors.New(`failed to download. Please check url`)
	}
	// detectedFileType := http.DetectContentType(bodyBytes)
	detectedFileType := mimetype.Detect(bodyBytes).String()

	if err != nil {
		fmt.Println(err.Error())
		return errors.New(`failed to download. Please check url`)
	}

	fName := uuid.New().String()
	// TODO : Use Content-Disposition header for fileName if available
	fileName := fName + getFileExtension(detectedFileType)
	tmpPath := filepath.Join(fPath, fileName)
	if getFileExtension(detectedFileType) == "" {
		logger.Error("Failed to get file extension")
		return errors.New(`failed to download. Please check url`)
	}
	newFile, err := os.Create(tmpPath)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(`failed to download. Please check url`)
	}

	defer newFile.Close()

	if _, err = newFile.Write(bodyBytes); err != nil {
		logger.Error(err.Error())
		return errors.New(`failed to download. Please check url`)
	}

	return
}

func getFileExtension(mimeType string) (extension string) {
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return
	}
	if len(ext) > 0 {
		return ext[len(ext)-1]
	}
	return ".png"
}
