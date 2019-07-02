# Crawl

[![Build Status](https://travis-ci.org/benmatselby/crawl.png?branch=master)](https://travis-ci.org/benmatselby/crawl)
[![codecov](https://codecov.io/gh/benmatselby/crawl/branch/master/graph/badge.svg)](https://codecov.io/gh/benmatselby/crawl)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/crawl?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/crawl)

The aim of this project was to write a simple web crawler. The crawler should be limited to one domain - so when you start with https://bbc.co.uk/, it would crawl all pages within bbc.co.uk, but not follow external links, for example to the Facebook and Twitter accounts. Given a URL, it should print a simple site map, showing the links between pages.

## Usage

```text

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
