all:
	cp -f ./prelude.v $(GOPATH)/bin/prelude.v
	go install

fmt:
	go fmt . ./...

test: fmt all
	go test ./...

testprelude: test
	vamos prelude_tests.v

run: fmt all
	vamos
