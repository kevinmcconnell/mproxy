build:
	CGO_ENABLED=0 go build -o bin/mproxy .

test:
	go test ./...

docker:
	docker build -t mproxy .
