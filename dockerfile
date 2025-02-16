FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

COPY .env .env

RUN mkdir -p /app/build \
    && go build -o /app/build/avito-shop ./cmd/avito-shop \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/app/build/avito-shop"]
