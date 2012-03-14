all:
	go install

fmt:
	go fmt . ./...

test: all fmt
	go test . ./...
