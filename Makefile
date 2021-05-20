export GO111MODULE = on

.PHONY: default test test-cover dev build

# for dev
dev:
	export GO_ENV=dev && fresh

build:
	go build -tags netgo -o web-screenshot

release:
	go mod tidy
