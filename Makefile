.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: gomod_tidy
gomod_tidy: ## Update Go mod dependencies
	 go mod tidy

.PHONY: gofmt
gofmt: ## Run Go formatting checks
	go fmt -x ./...

.PHONY: build
build: ## Build executable without extension (Unix)
	go build -o alizer

.PHONY: buildWin
buildWin: ## Build executable with Windows .exe extension
	go build -o alizer.exe

.PHONY: test
test: ## Run unit tests with coverage
	go test -coverprofile cover.out -v ./...

.PHONY: check_registry
check_registry: ## Run registry checks
	./test/check_registry/check_registry.sh

.PHONY: gosec_install
gosec_install: ## Install gosec utility
	go install github.com/securego/gosec/v2/cmd/gosec@v2.14.0

.PHONY: gosec
gosec: ## Run go security checks
	./run_gosec.sh
