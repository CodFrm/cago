
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
	cago swag

lint: check-golangci-lint
	golangci-lint run

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
