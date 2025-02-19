GOBIN ?= $(shell go env GOPATH)/bin
GOFMT:=gofumpt

default: build

dep: # Download required dependencies
	go mod tidy

# build: fmtcheck
# 	go build -o $(GOBIN)/okta-cli cmd/okta-cli/main.go

build:
	go build -o $(GOBIN)/okta-cli-client . 

install:
	go install .

generate-cmd:
	go run ./cmdTools

generate-sdk:
	openapi-generator-cli version-manager set 7.0.1 \
	&& npx @openapitools/openapi-generator-cli generate -c ./.generator/config.yaml -i ./template.yaml \
	&& cd sdk \
	&& rm -rf git_push.sh go.mod go.sum

# fmtcheck:
# 	@gofumpt -d -l .

fmt: tools # Format the code
	@$(GOFMT) -l -w .

test:
	go test -race -v $(TEST) || exit 1

tools:
	@which $(ERRCHECK) || go install github.com/kisielk/errcheck@latest
	@which $(GOCONST) || go install github.com/jgautheron/goconst/cmd/goconst@latest
	@which $(GOCYCLO) || go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@which $(GOFMT) || go install mvdan.cc/gofumpt@latest
	@which $(GOLINT) || go install golang.org/x/lint/golint@latest
	@which $(SHADOW) || go mod download golang.org/x/tools
	@which $(STATICCHECK) || go install honnef.co/go/tools/cmd/staticcheck@latest

tools-update:
	@go install github.com/kisielk/errcheck@latest
	@go install github.com/jgautheron/goconst/cmd/goconst@latest
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@go install mvdan.cc/gofumpt@latest
	@go install golang.org/x/lint/golint@latest
	@go mod download golang.org/x/tools
	@go install honnef.co/go/tools/cmd/staticcheck@latest
