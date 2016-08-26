all: build

test:
	go test -v -cover ./...

lint:
	go vet ./...
	golint -set_exit_status ./...

deps:
	go get -d -t ./...
	go get github.com/golang/lint/golint

build: test lint clean
	GOOS=darwin GOARCH=amd64 go build -o build/chaosmonkey_darwin_amd64
	GOOS=linux  GOARCH=amd64 go build -o build/chaosmonkey_linux_amd64
	cd build && shasum -a256 chaosmonkey_* > SHA256SUMS

clean:
	$(RM) -r build

.PHONY: build
