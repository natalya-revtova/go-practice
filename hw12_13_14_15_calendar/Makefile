BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GOOPENAPIV2_VERSION?=v2@v2.17.0
GOOGLEAPIS_VERSION?=v1.16.0

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DB_USER = $(CALENDAR_DB_USER)
DB_PSSWD = $(CALENDAR_DB_PSSWD)
DB_HOST := $(CALENDAR_DB_HOST)
DB_PORT := 5432
DB_NAME := calendar

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.52.2

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

migrate:
	 cd migrations && goose postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PSSWD) dbname=$(DB_NAME)" up

.PHONY: build run build-img run-img version test lint