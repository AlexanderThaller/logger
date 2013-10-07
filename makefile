format:
	gofmt -s=true -w=true *.go

test:
	go test

bench:
	go test -test.bench .
