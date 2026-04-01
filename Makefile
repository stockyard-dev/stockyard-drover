build:
	CGO_ENABLED=0 go build -o drover ./cmd/drover/

run: build
	./drover

test:
	go test ./...

clean:
	rm -f drover

.PHONY: build run test clean
