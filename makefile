format:
	gofmt -s=true -w=true logger.go logger_test.go

test:
	go test
