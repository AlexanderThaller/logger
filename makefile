format:
	gofmt -s=true -w=true logger.go logger_test.go

test:
	go test

bench:
	go test -test.bench Log
	go test -test.bench get
	go test -test.bench print
	go test -test.bench format
