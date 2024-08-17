package toolkit

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
)

func GetTitleByUrl(rawUrl string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	resp, err := http.Get(parsedUrl.String())
	if err != nil {
		return "", fmt.Errorf("error while fetching URL: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error while closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error while parsing document: %v", err)
		}
		return doc.Find("title").Text(), nil
	}
	return "Error while fetching title.", nil
}
