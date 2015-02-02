all:
	cp -f ./prelude.v ../../bin/prelude.v
	go install

fmt:
	go fmt . ./...

test: fmt all
	go test ./...

run: fmt all
	vamos
