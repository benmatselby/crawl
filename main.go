package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	"golang.org/x/net/html"
)

// Logger will allow us to use a --verbose flag
type Logger struct {
	Verbose bool
	Writer  io.Writer
}

// Write checks to see if we want verbosity and handles the log accordingly
func (l Logger) Write(data []byte) (n int, err error) {
	if l.Verbose {
		return fmt.Fprint(l.Writer, string(data))
	}

	return 0, nil
}

// Page houses information about the page
type Page struct {
	URL       string
	HasErrors bool
	Links     map[string]int
}

// SiteMap details what the JSON will look like
type SiteMap struct {
	Pages   []Page `json:"pages"`
	Depth   int
	WasMore bool
}

// visited stores the urls we have crawled
var visited = struct {
	urls map[string]bool
	sync.Mutex
}{urls: make(map[string]bool)}

// sitemap keeps a record
var sitemap = SiteMap{}

// desiredDepth determines how deep into the crawl we go before we stop
var desiredDepth int

// main is the entry point to the application
func main() {
	var verbosity bool
	flag.BoolVar(&verbosity, "verbose", false, "whether we want to output all the crawling information")
	flag.IntVar(&desiredDepth, "depth", 2, "how deep into the crawl before we stop")
	flag.Parse()

	l := Logger{Verbose: verbosity, Writer: os.Stdout}
	err := Run(flag.Args(), l)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	json, err := json.Marshal(sitemap)
	fmt.Println(string(json))
}

// Run is the main execution of the application
func Run(args []string, w io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("no URL specified")
	}

	sitemap.Depth = desiredDepth

	target, err := url.ParseRequestURI(args[0])
	if err != nil {
		return fmt.Errorf("invalid URL")
	}
	CrawlPage(target.String(), 0, w)

	return nil
}

// CrawlPage will scan a single page and then generate more go routines to carry on down the rabbit hole
func CrawlPage(pageURL string, level int, w io.Writer) {
	if level > desiredDepth {
		fmt.Fprintf(w, "not fetching %s as we are at depth of %v\n", pageURL, level)
		sitemap.WasMore = true
		return
	}
	fmt.Fprintf(w, "fetching %s\n", pageURL)

	page := Page{
		URL: pageURL,
	}

	response, err := http.Get(pageURL)
	if err != nil {
		page.HasErrors = true
		fmt.Fprintf(w, "error fetching url %s: %s", pageURL, err)
		return
	}

	visited.Lock()
	if visited.urls[pageURL] == true {
		fmt.Fprintf(w, "already crawled %s\n", pageURL)
		visited.Unlock()
		return
	}
	visited.urls[pageURL] = true
	visited.Unlock()

	done := make(chan bool)

	targetURL, err := url.ParseRequestURI(pageURL)
	if err != nil {
		fmt.Fprintf(w, "error parsing URL %s: %s", pageURL, err)
	}

	urls := GetURLs(targetURL, response.Body)
	page.Links = urls
	sitemap.Pages = append(sitemap.Pages, page)

	for url := range urls {
		go func(url string, level int, w io.Writer) {
			CrawlPage(url, level, w)
			done <- true
		}(url, level+1, w)
	}

	for url := range urls {
		fmt.Fprintf(w, "waiting %s\n", url)
		<-done
	}
}

// GetURLs is responsible for parsing the HTML document, and returning usable URLs
func GetURLs(currentURL *url.URL, body io.ReadCloser) map[string]int {
	var urls = map[string]int{}
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		switch tokenType {
		case html.ErrorToken:
			return urls
		case html.StartTagToken:
			if token.DataAtom.String() == "a" {
				for _, value := range token.Attr {
					if value.Key == "href" {
						foundURL, err := url.Parse(value.Val)
						if err != nil {
							continue
						}
						visited.Lock()
						if foundURL.Host == currentURL.Host {
							urls[foundURL.String()] = urls[foundURL.String()] + 1
						}
						visited.Unlock()
					}
				}
			}
		}
	}
}
