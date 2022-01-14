BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
GOOPENAPIV2_VERSION?=v2@v2.7.2
GOOGLEAPIS_VERSION?=v1.16.0

VERSION?=$(shell git log --date=short --pretty=format:'%ad-%h' -n 1 | sed -E 's/([0-9]{4})\-([0-9]{2})\-([0-9]{2})/\1.\2.\3/g')

LDFLAGS := -X github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config.VERSION=$(VERSION)

generate:
	rm -rf proto/pb
	mkdir -p proto/pb
	protoc -I. \
			-I$(GOPATH)/pkg/mod \
			-I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GOOGLEAPIS_VERSION)/third_party/googleapis \
			-I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/$(GOOPENAPIV2_VERSION) \
			--go_out=proto/pb \
			--go-grpc_out=proto/pb \
			--grpc-gateway_out=proto/pb \
			proto/*.proto
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

.PHONY: up down generate migration-up migration-down build run build-img run-img test install-lint-deps lint
