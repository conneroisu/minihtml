// Package main is a simple command line tool that takes a URL as input and
// outputs a cleaned up HTML file.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// main is the entry point for the program.
func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// run takes a URL as input and outputs a cleaned up HTML file.
func run(ctx context.Context, args []string) error {
	// take the url from the command line
	url := args[1]
	fmt.Println("url:", url)
	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	sanitized := colly.SanitizeFileName(url)
	// create a new file
	f, err := os.Create(sanitized)
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := clean(body)
	// write the body to the file
	_, err = f.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

// clean removes all the non-essential elements from the HTML document.
func clean(body []byte) (string, error) {
	// get a goquery document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	doc.Find("script").Remove()
	doc.Find("style").Remove()
	doc.Find("head").Remove()
	doc.Find("meta").Remove()
	doc.Find("link").Remove()
	doc.Find("img").Remove()
	doc.Find("form").Remove()
	doc.Find("input").Remove()
	doc.Find("button").Remove()
	return doc.Html()
}
