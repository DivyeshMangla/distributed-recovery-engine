BINARY := dre
CMD := ./cmd/dre

.PHONY: build run test tidy clean node

build:
	go build -o bin/$(BINARY) $(CMD)

run:
	go run $(CMD)

node:
	go run $(CMD) --id=$(ID) --addr=$(ADDR) --seed=$(SEED)

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -rf bin/
