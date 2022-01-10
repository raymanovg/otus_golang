BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

VERSION?=$(shell git log --date=short --pretty=format:'%ad-%h' -n 1 | sed -E 's/([0-9]{4})\-([0-9]{2})\-([0-9]{2})/\1.\2.\3/g')

LDFLAGS := -X github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config.VERSION=$(VERSION)

generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb

	protoc \
		--proto_path=api/ \
		--go_out=internal/server/grpc/pb \
		--go-grpc_out=internal/server/grpc/pb \
		api/*.proto

up:
	docker-compose -f deployments/docker-compose.yaml up -d calendar

down:
	docker-compose -f deployments/docker-compose.yaml down

rebuild: build-img up

migration-up:
	docker-compose -f deployments/docker-compose.yaml up migration-up

migration-down:
	docker-compose -f deployments/docker-compose.yaml up migration-down

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
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
