# buildtime
FROM golang:1.24 AS build

WORKDIR /app

# COPY go.mod go.sum ./  when the project eventually has dependencies
COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o dre ./cmd/dre

# runtime
FROM scratch

WORKDIR /app

COPY --from=build /app/dre /app/dre

EXPOSE 8080

ENTRYPOINT ["/app/dre"]
