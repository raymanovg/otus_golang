BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

VERSION?=$(shell git log --date=short --pretty=format:'%ad-%h' -n 1 | sed -E 's/([0-9]{4})\-([0-9]{2})\-([0-9]{2})/\1.\2.\3/g')

LDFLAGS := -X github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config.VERSION=$(VERSION)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) --config=./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
