.PHONY: \
build \
test test-verbose cover-report \
run \
logs logs-follow logs-app logs-app-follow \
clean stop \
rsa \
mockgen

BINARY_NAME=moneybags
COMPOSE_SERVICE_NAME=moneybags

# GO_SOURCES=$(shell find . -path "./mock" -prune -o -name "*.go") go.mod go.sum

build: clean # $(GO_SOURCES)
	GO111MODULE=on \
	GOOS=linux \
	GOARCH=amd64 \
	CGO_ENABLED=0 \
	go build \
	-o $(BINARY_NAME) cmd/main.go

test: mockgen
	go test ./... -cover

test-verbose: mockgen
	go test -v ./... -cover

cover-report: mockgen
	-go test -coverprofile=coverage.out ./...
	-go tool cover -html=coverage.out -o coverage.html

run: stop build rsa
	docker-compose build
	docker-compose up -d

logs: 
	docker-compose logs

logs-follow:
	docker-compose logs -f

logs-app:
	docker-compose logs $(COMPOSE_SERVICE_NAME)

logs-app-follow:
	docker-compose logs $(COMPOSE_SERVICE_NAME)

clean:
	-rm $(BINARY_NAME)
	-rm -r ./secrets

clean-mocks:
	-rm -r ./mock

stop:
	-docker-compose down

rsa:
	-mkdir ./secrets
	openssl genrsa -out ./secrets/private.pem 3072
	openssl rsa -in ./secrets/private.pem -pubout -out ./secrets/public.pem

mockgen: clean-mocks
	-go generate ./...