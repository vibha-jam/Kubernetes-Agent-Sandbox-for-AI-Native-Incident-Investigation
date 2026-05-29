FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o controller \
    ./cmd/controller/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/controller .

CMD ["./controller"]