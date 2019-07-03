package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

// Page houses information about the page
type Page struct {
	URL   string
	Links []string
}

// SiteMap details what the JSON will look like
type SiteMap struct {
	Pages []Page `json:"pages"`
}

// main is the entry point to the application
func main() {
	flag.Parse()
	sitemap, err := Run(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	json, err := json.Marshal(sitemap)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(json))
}

// Run is the main execution of the application
func Run(args []string) (*SiteMap, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("no URL specified")
	}

	target, err := url.ParseRequestURI(args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid URL")
	}

	var pages []Page

	found, err := CrawlPage(target)
	if err != nil {
		// @todo
	}

	pages = append(pages, Page{URL: target.String(), Links: found})
	sitemap := &SiteMap{
		Pages: pages,
	}

	return sitemap, nil
}

// CrawlPage will scan a single page and return the URLs it finds if they match the target
func CrawlPage(URL *url.URL) ([]string, error) {
	response, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	var urls []string
	tokenizer := html.NewTokenizer(response.Body)
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		switch tokenType {
		case html.ErrorToken:
			return urls, nil
		case html.StartTagToken:
			if token.DataAtom.String() == "a" {
				for _, value := range token.Attr {
					if value.Key == "href" {
						pageURL, err := url.Parse(value.Val)
						if err != nil {
							continue
						}

						if pageURL.Host == URL.Host {
							urls = append(urls, pageURL.String())
						}
					}
				}
			}
		}
	}
}
