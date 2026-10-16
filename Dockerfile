# Multi-stage build for orchid-auth-starter (Orchid Starter Clean Architecture)
FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/api ./cmd/api

# Final minimal runtime image
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/api /app/bin/api
COPY .env.example /app/.env.example

EXPOSE 8080

USER nobody:nobody

ENTRYPOINT ["/app/bin/api"]
