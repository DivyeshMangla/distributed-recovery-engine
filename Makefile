BINARY := dre
CMD := ./cmd/dre

.PHONY: build run tidy test clean

build:
	go build -o bin/$(BINARY) $(CMD)

run:
	go run $(CMD)

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -rf bin/
