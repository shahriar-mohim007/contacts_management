FROM golang:1.23.1 as builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o api cmd/api/main.go

FROM alpine:latest



WORKDIR /app
COPY --from=builder /app/api /usr/local/bin/api
CMD ["/usr/local/bin/api"]
