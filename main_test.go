package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

func TestRunCanParseFlags(t *testing.T) {
	tt := []struct {
		name     string
		args     []string
		expected error
	}{
		{name: "no arguments passed in", args: []string{}, expected: errors.New("no URL specified")},
		{name: "bad url", args: []string{"flim flam"}, expected: errors.New("invalid URL")},
		{name: "valid url", args: []string{"https://bbc.co.uk"}, expected: nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := Run(tc.args)
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
		response string
		goodURLs []string
		badURLs  []string
		expected []string
	}{
		{
			name: "a single URL",
			URL:  "/single-url",
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
			expected: []string{"http://GOOD_URL/blog"},
		},
		{
			name: "a single good URL with other urls",
			URL:  "/single-url-with-other-urls",
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
			expected: []string{"http://GOOD_URL/blog"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			re := regexp.MustCompile(`GOOD_URL`)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, re.ReplaceAllString(tc.response, r.Host))
			}))
			defer ts.Close()

			URL, err := url.ParseRequestURI(ts.URL + tc.URL)
			if err != nil {
				t.Fatalf("was not expected error, got '%s'", err)
			}

			urls, err := CrawlPage(URL)
			if err != nil {
				t.Fatalf("was not expected error, got '%s'", err)
			}

			testServerURL, err := url.ParseRequestURI(ts.URL)
			if err != nil {
				t.Fatalf("was not expected error, got '%s'", err)
			}

			var expected []string
			for _, url := range tc.expected {
				expected = append(expected, re.ReplaceAllString(url, testServerURL.Host))
			}

			if fmt.Sprintf("%s", urls) != fmt.Sprintf("%s", expected) {
				t.Fatalf("expected '%s', got '%s'", tc.expected, urls)
			}
		})
	}
}
