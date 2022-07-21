
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
	ETCD_ADDR=192.168.110.230:2379 FC_APP=fc-cspm NAMESPACE=dev go test -v ./...

coverage.out cover:
	ETCD_ADDR=192.168.110.230:2379 FC_APP=fc-cspm NAMESPACE=dev go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

html-cover: coverage.out
	go tool cover -html=coverage.out
	go tool cover -func=coverage.out

generate: check-mockgen swagger
	go generate ./... -x
