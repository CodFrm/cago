
check-swag:
ifneq ($(which swag),)
	go get -u github.com/swaggo/swag/cmd/swag
endif

check-mockgen:
ifneq ($(which mockgen),)
	go install github.com/golang/mock/mockgen
endif

swagger: check-swag
	swag init

lint:
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
