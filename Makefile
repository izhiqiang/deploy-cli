GIT_COMMIT:=$(shell git rev-list -1 HEAD)
LDFLAGS:=-X main.GitCommit=${GIT_COMMIT}

build:
	go build -ldflags "$(LDFLAGS)" .

check:
	golangci-lint run ./...