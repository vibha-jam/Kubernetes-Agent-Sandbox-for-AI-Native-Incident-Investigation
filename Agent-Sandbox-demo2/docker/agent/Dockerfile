FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o ai-agent \
    ./cmd/agent/main.go


FROM ollama/ollama:latest

WORKDIR /app

COPY --from=builder /app/ai-agent .

EXPOSE 8080
EXPOSE 11434

COPY start.sh .

RUN chmod +x start.sh

ENTRYPOINT []

CMD ["./start.sh"]