
check-cago:
ifneq ($(which cago),)
	go install github.com/codfrm/cago
endif

check-mockgen:
ifneq ($(which mockgen),)
	go install github.com/golang/mock/mockgen
endif

check-golangci-lint:
ifneq ($(which golangci-lint),)
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
endif

swagger: check-cago
	cago gen swag

lint: check-golangci-lint
	golangci-lint run

lint-fix: check-golangci-lint
	golangci-lint run --fix

test: lint
	go test -v ./...

coverage.out cover:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

html-cover: coverage.out
	go tool cover -html=coverage.out
	go tool cover -func=coverage.out

generate: check-mockgen swagger
	go generate ./... -x

GOOS=linux
GOARCH=amd64
APP_NAME=simple
APP_VERSION=1.0.0

SUFFIX=
ifeq ($(GOOS),windows)
	SUFFIX=.exe
endif

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(APP_NAME)_v$(APP_VERSION)$(SUFFIX) ./cmd/app

