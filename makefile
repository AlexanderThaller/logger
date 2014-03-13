default:
	make format
	make test

format:
	gofmt -s=true -w=true *.go
	goimports -w=true *.go
# golint *.go */*/*.go

test:
	go test

bench:
	go test -test.bench .
