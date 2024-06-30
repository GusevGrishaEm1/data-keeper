FROM golang:1.22.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/datakeeper/main.go

FROM docker:20.10.24-dind AS tester

ENTRYPOINT ["sh", "-c", "dockerd-entrypoint.sh & sleep 3 && go test ./..."]

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["./server"]