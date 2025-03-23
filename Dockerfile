FROM golang:1.24

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/main.go

CMD ["/app/main"]
