export GO111MODULE = on

.PHONY: default test test-cover dev build

# for dev
dev:
	export GO_ENV=dev && air 

build:
	go build -tags netgo -o web-screenshot

release:
	go mod tidy
