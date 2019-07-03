# Crawl

[![Build Status](https://travis-ci.org/benmatselby/crawl.png?branch=master)](https://travis-ci.org/benmatselby/crawl)
[![codecov](https://codecov.io/gh/benmatselby/crawl/branch/master/graph/badge.svg)](https://codecov.io/gh/benmatselby/crawl)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/crawl)](https://goreportcard.com/report/github.com/benmatselby/crawl)

The aim of this project was to write a simple web crawler. The crawler should be limited to one domain - so when you start with [https://bbc.co.uk/](https://bbc.co.uk/), it would crawl all pages within [bbc.co.uk](https://bbc.co.uk/), but not follow external links, for example to the Facebook and Twitter accounts. Given a URL, it should print a simple site map, showing the links between pages.

## Usage

```text
crawl https://bbc.co.uk
```

## Requirements

- [Go version 1.12+](https://golang.org/dl/)

## Installation via Git

```bash
git clone git@github.com:benmatselby/crawl.git
cd crawl
make all
./crawl
```

You can also install into your `$GOPATH/bin` by `go install`

## Future

- [ ] Define site map writers using the strategy pattern (e.g. have an xml, json writer).
- [ ] Have throttling in the system.
- [ ] Observe the `robots.txt`.
- [ ] Provide a mechanism to highlight broken links.
- [ ] Provide a mechanism to run for assets like JavaScript, CSS etc.
- [ ] Cater for a site that may use relative paths such as "/" so the URL works on any domain.
