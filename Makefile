test: lint vet

vet:
	go vet ./...

lint:
	golint ./...
