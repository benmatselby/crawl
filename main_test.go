package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

func TestLogger(t *testing.T) {
	tt := []struct {
		name      string
		verbosity bool
		msg       string
		expected  string
	}{
		{name: "verbose mode", verbosity: true, msg: "My logged string", expected: "My logged string"},
		{name: "verbose mode off", verbosity: false, msg: "My logged string", expected: ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			writer := bufio.NewWriter(&b)

			logger := Logger{
				Verbose: tc.verbosity,
				Writer:  writer,
			}

			fmt.Fprintf(logger, tc.msg)
			writer.Flush()

			if b.String() != tc.expected {
				t.Fatalf("expected '%s'; got '%s'", tc.expected, b.String())
			}
		})
	}
}

func TestRunCanParseFlags(t *testing.T) {
	tt := []struct {
		name     string
		args     []string
		expected error
	}{
		{name: "no arguments passed in", args: []string{}, expected: errors.New("no URL specified")},
		{name: "bad url", args: []string{"flim flam"}, expected: errors.New("invalid URL")},
		{name: "valid url", args: []string{"http://GOOD_URL"}, expected: nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sitemap.Pages = nil
			sitemap = SiteMap{}
			visited.urls = nil
			visited.urls = make(map[string]bool)
			desiredDepth = 0

			re := regexp.MustCompile(`GOOD_URL`)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "")
			}))
			defer ts.Close()

			testServerURL, err := url.ParseRequestURI(ts.URL)
			if err != nil {
				t.Fatalf("was not expected error, got '%s'", err)
			}

			if len(tc.args) == 1 {
				tc.args[0] = re.ReplaceAllString(tc.args[0], testServerURL.Host)
			}

			var b bytes.Buffer
			w := bufio.NewWriter(&b)

			err = Run(tc.args, w)
			if tc.expected == nil && err != nil {
				t.Fatalf("did not expect error, got %v", err)
			}

			if tc.expected != nil {
				if err.Error() != tc.expected.Error() {
					t.Fatalf("expected %v, got %v", tc.expected, err)
				}
			}
		})
	}
}

func TestCrawlPage(t *testing.T) {
	tt := []struct {
		name     string
		URL      string
		depth    int
		response string
		goodURLs []string
		badURLs  []string
		expected SiteMap
	}{
		{
			name:  "a single URL",
			URL:   "/single-url",
			depth: 10,
			response: `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<title>Document</title>
			</head>
			<body>
			<p>A paragraph</p>
			<a href="http://GOOD_URL/blog">a link</a>
			</body>
			</html>`,
			expected: SiteMap{
				WasMore: false,
				Pages: []Page{
					Page{
						URL:   "http://GOOD_URL/single-url",
						Links: map[string]int{"http://GOOD_URL/blog": 1},
					},
					Page{
						URL:   "http://GOOD_URL/blog",
						Links: map[string]int{"http://GOOD_URL/blog": 1},
					},
				},
			},
		},
		{
			name:  "a single good URL with other urls",
			URL:   "/single-url-with-other-urls",
			depth: 10,
			response: `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<title>Document</title>
			</head>
			<body>
			<p>A paragraph</p>
		  <a href="http://GOOD_URL/blog">a link</a>
		  <a href="https://twitter.com/">a link</a>
		  <a href="https://github.com/">a link</a>
			</body>
			</html>`,
			expected: SiteMap{
				WasMore: false,
				Pages: []Page{
					Page{
						URL:   "http://GOOD_URL/single-url-with-other-urls",
						Links: map[string]int{"http://GOOD_URL/blog": 1},
					},
					Page{
						URL:   "http://GOOD_URL/blog",
						Links: map[string]int{"http://GOOD_URL/blog": 1},
					},
				},
			},
		},
		{
			name:  "multiple good urls, no bad urls",
			URL:   "/multiple-urls",
			depth: 10,
			response: `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<title>Document</title>
			</head>
			<body>
			<p>A paragraph</p>
		  <a href="http://GOOD_URL/blog">a link</a>
		  <a href="http://GOOD_URL/about-us">a link</a>
			</body>
			</html>`,
			expected: SiteMap{
				WasMore: false,
				Pages: []Page{
					Page{
						URL:   "http://GOOD_URL/multiple-urls",
						Links: map[string]int{"http://GOOD_URL/blog": 1, "http://GOOD_URL/about-us": 1},
					},
					Page{
						URL:   "http://GOOD_URL/blog",
						Links: map[string]int{"http://GOOD_URL/blog": 1, "http://GOOD_URL/about-us": 1},
					},
					Page{
						URL:   "http://GOOD_URL/about-us",
						Links: map[string]int{"http://GOOD_URL/blog": 1, "http://GOOD_URL/about-us": 1},
					},
				},
			},
		},
		{
			name:  "multiple of the same url",
			URL:   "/multiple-of-same-url",
			depth: 10,
			response: `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<title>Document</title>
			</head>
			<body>
			<p>A paragraph</p>
		  <a href="http://GOOD_URL/blog">a link</a>
		  <a href="http://GOOD_URL/blog">a link</a>
			</body>
			</html>`,
			expected: SiteMap{
				WasMore: false,
				Pages: []Page{
					Page{
						URL:   "http://GOOD_URL/multiple-of-same-url",
						Links: map[string]int{"http://GOOD_URL/blog": 1},
					},
					Page{
						URL:   "http://GOOD_URL/blog",
						Links: map[string]int{"http://GOOD_URL/blog": 1},
					},
				},
			},
		},
		{
			name:  "no urls",
			URL:   "/no-urls",
			depth: 10,
			response: `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<title>Document</title>
			</head>
			<body>
			<p>A paragraph</p>
			</body>
			</html>`,
			expected: SiteMap{
				WasMore: false,
				Pages: []Page{
					Page{
						URL:   "http://GOOD_URL/no-urls",
						Links: map[string]int{},
					},
				},
			},
		},
		{
			name:  "more urls than we want to track",
			URL:   "/more-urls-than-we-care-about",
			depth: 0,
			response: `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<title>Document</title>
			</head>
			<body>
			<p>A paragraph</p>
		  <a href="http://GOOD_URL/blog">a link</a>
		  <a href="http://GOOD_URL/blog">a link</a>
			</body>
			</html>`,
			expected: SiteMap{
				WasMore: true,
				Pages: []Page{
					Page{
						URL:   "http://GOOD_URL/more-urls-than-we-care-about",
						Links: map[string]int{"http://GOOD_URL/blog": 1},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sitemap.Pages = nil
			sitemap = SiteMap{}
			visited.urls = nil
			visited.urls = make(map[string]bool)
			desiredDepth = tc.depth

			re := regexp.MustCompile(`GOOD_URL`)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, re.ReplaceAllString(tc.response, r.Host))
			}))
			defer ts.Close()

			URL, err := url.ParseRequestURI(ts.URL + tc.URL)
			if err != nil {
				t.Fatalf("was not expected error, got '%s'", err)
			}

			var b bytes.Buffer
			w := bufio.NewWriter(&b)

			CrawlPage(URL.String(), 0, w)
			if err != nil {
				t.Fatalf("was not expected error, got '%s'", err)
			}

			testServerURL, err := url.ParseRequestURI(ts.URL)
			if err != nil {
				t.Fatalf("was not expected error, got '%s'", err)
			}

			if len(sitemap.Pages) != len(tc.expected.Pages) {
				t.Fatalf("expected %v pages, got %v", len(sitemap.Pages), len(tc.expected.Pages))
			}

			if sitemap.WasMore != tc.expected.WasMore {
				t.Fatalf("expected WasMore to be '%v', got '%v'", tc.expected.WasMore, sitemap.WasMore)
			}

			for _, expectedPage := range tc.expected.Pages {
				expectedURL := re.ReplaceAllString(expectedPage.URL, testServerURL.Host)
				foundPage := false
				for _, p := range sitemap.Pages {
					if expectedURL == p.URL {
						foundPage = true

						if len(expectedPage.Links) != len(p.Links) {
							t.Fatalf("expected page %s to have %v links, got %v", expectedPage.URL, len(expectedPage.Links), len(p.Links))
						}

						for expectedLink := range expectedPage.Links {
							expectedLink := re.ReplaceAllString(expectedLink, testServerURL.Host)
							foundLink := false
							for l := range p.Links {
								if expectedLink == l {
									foundLink = true
								}
							}

							if !foundLink {
								t.Fatalf("expected to find link %s in page %s, but did not", expectedLink, expectedPage.URL)
							}
						}
					}
				}

				if !foundPage {
					t.Fatalf("expected to find page %s, but did not", expectedURL)
				}
			}
		})
	}
}
