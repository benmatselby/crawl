# Crawl

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=crawl&metric=alert_status)](https://sonarcloud.io/dashboard?id=crawl)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/crawl)](https://goreportcard.com/report/github.com/benmatselby/crawl)

The aim of this project was to write a simple web crawler. The crawler should be limited to one domain - so when you start with [https://benmatselby.dev/](https://benmatselby.dev/), it would crawl all pages within [benmatselby.dev](https://benmatselby.dev/), but not follow external links, for example to the Facebook and Twitter accounts. Given a URL, it should print a simple site map, showing the links between pages.

## Usage

```shell
Usage of crawl:
  -depth int
      how deep into the crawl before we stop (default 2)
  -verbose
      whether we want to output all the crawling information

crawl https://benmatselby.dev
```

This will return a JSON string for the output. You can then pipe that into something like [jq](https://stedolan.github.io/jq/) to get some nice formatting, or interrogation options:

**URLs crawled**:

```shell
./crawl https://benmatselby.dev | jq '.pages[].URL'
"https://benmatselby.dev"
"https://benmatselby.dev/"
"https://benmatselby.dev/post/sugarcrm-deployment-process/"
"https://benmatselby.dev/post/"
"https://benmatselby.dev/post/joining-a-new-engineering-team/"
"https://benmatselby.dev/post/software-engineering-team-structure/"
"https://benmatselby.dev/post/pa11y-accessibility-ci/"
"https://benmatselby.dev/post/development-environments/"
"https://benmatselby.dev/post/pipelines/"
"https://benmatselby.dev/post/communication/"
"https://benmatselby.dev/post/feature-toggling/"
"https://benmatselby.dev/post/onboarding/"
"https://benmatselby.dev/post/technology-radar/"
"https://benmatselby.dev/post/communication-tools/"
"https://benmatselby.dev/post/squashing-commits/"
"https://benmatselby.dev/post/why-teams-important/"
```

**Count of URLs crawled**:

```shell
./crawl https://benmatselby.dev | jq '.pages | length'
16
```

## Requirements

- [Go version 1.12+](https://golang.org/dl/)

## Installation via Git

The main way to install this application, is to clone the git repo, and build. This will require you to have the Go runtime installed on your machine.

```shell
git clone git@github.com:benmatselby/crawl.git
cd crawl
make all
./crawl
```

Once built, you can also install the binary into your `$GOPATH/bin` by `go install`. This will mean that `crawl` will be globally available in your system.

## Installation via Docker

Whilst this is not the recommended way to run the application, as there is a slight performance overhead in running it in a container, you can do so.

```shell
git clone git@github.com:benmatselby/crawl.git
cd crawl
make build-docker
docker run benmatselby/crawl https://benmatselby.dev
```

## Future

- [ ] Define site map writers using the strategy pattern (e.g. have an xml, json writer).
- [ ] Have throttling in the system.
- [ ] Observe the `robots.txt`.
- [ ] Provide a mechanism to highlight broken links.
- [ ] Provide a mechanism to run for assets like JavaScript, CSS etc.
- [ ] Cater for a site that may use relative paths such as "/" so the URL works on any domain.
