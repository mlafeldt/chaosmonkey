test:
	go test -v -cover ./...

lint:
	go vet ./...
	golint -set_exit_status ./...

deps:
	go get \
		github.com/aws/aws-sdk-go/aws/... \
		github.com/aws/aws-sdk-go/service/... \
		github.com/golang/lint/golint \
		github.com/ryanuber/columnize

build: test lint
	GOOS=darwin GOARCH=amd64 go build -o build/chaosmonkey_darwin_amd64
	GOOS=linux  GOARCH=amd64 go build -o build/chaosmonkey_linux_amd64

.PHONY: build
