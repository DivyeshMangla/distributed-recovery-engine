FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o dre ./cmd/dre

EXPOSE 8080

ENTRYPOINT ["./dre"]