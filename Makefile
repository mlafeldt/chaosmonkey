PACKAGES := $(shell go list ./... | grep -v vendor)

test:
	go test -v -cover $(PACKAGES)

lint:
	go vet $(PACKAGES)
	echo $(PACKAGES) | xargs -n1 golint -set_exit_status
