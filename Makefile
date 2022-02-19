.PHONY: \
build \
test test-verbose\
run \
logs logs-follow \
clean stop

BINARY_NAME=moneybags

GO_SOURCES=$(shell find . -name "*.go") go.mod go.sum

build: clean $(GO_SOURCES)
	GO111MODULE=on \
	GOOS=linux \
	GOARCH=amd64 \
	CGO_ENABLED=0 \
	go build \
	-o $(BINARY_NAME) cmd/main.go

test:
	go test ./...

test-verbose:
	go test -v ./...

run: stop build
	docker-compose build
	docker-compose up -d

logs:
	docker-compose logs

logs-follow:
	docker-compose logs -f

clean:
	-rm $(BINARY_NAME)

stop:
	-docker-compose down