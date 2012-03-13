all:
	go install

test: all
	go test . ./...
