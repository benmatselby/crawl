APP_NAME := crawl
DOCKER_OWNER := benmatselby

.PHONY: explain
explain:
	### Welcome
	#
  #   ______ .______          ___   ____    __    ____  __
  #  /      ||   _  \        /   \  \   \  /  \  /   / |  |
  # |  ,----'|  |_)  |      /  ^  \  \   \/    \/   /  |  |
  # |  |     |      /      /  /_\  \  \            /   |  |
  # |  `----.|  |\  \----./  _____  \  \    /\    /    |  `----.
  #  \______|| _| `._____/__/     \__\  \__/  \__/     |_______|
  #
	#
	### Targets
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: vet
vet: ## Vet the code
	go vet -v ./...

.PHONY: lint
lint: ## Lint the code
	golint -set_exit_status $(shell go list ./... | grep -v vendor)

.PHONY: build
build: ## Build the application
	go build .

.PHONY: build-static
build-static: ## Build the application
	CGO_ENABLED=0 go build -ldflags "-extldflags -static" -o $(APP_NAME) .

.PHONY: build-docker
build-docker: ## Build the docker image
	docker build -t ${DOCKER_OWNER}/${APP_NAME} .

.PHONY: test
test: ## Run the unit tests
	go test ./... -coverprofile=coverage.out

.PHONY: test-cov
test-cov: test ## Run the unit tests with coverage
	go tool cover -html=coverage.out
