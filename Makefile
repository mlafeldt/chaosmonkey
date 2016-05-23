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
