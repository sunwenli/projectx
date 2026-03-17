build:
	go build -o ./bin/projectx.exe

run: build
	./bin/projectx.exe

test:
	go test -v ./...