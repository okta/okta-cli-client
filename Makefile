GOBIN ?= $(shell go env GOPATH)/bin
GOFMT:=gofumpt

default: build

# build: fmtcheck
# 	go build -o $(GOBIN)/okta-cli cmd/okta-cli/main.go

build:
	go build -o $(GOBIN)/okta-cli okta-cli/main.go

generate-cmd:
	go run ./cmdTools

generate-sdk:
	openapi-generator-cli version-manager set 7.0.1 \
	&& npx @openapitools/openapi-generator-cli generate -c ./.generator/config.yaml -i ./template.yaml \
	&& cd sdk \
	&& rm -rf git_push.sh go.mod go.sum

# fmtcheck:
# 	@gofumpt -d -l .