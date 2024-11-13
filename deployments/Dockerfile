FROM golang:1.23.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o woodpecker ./cmd/woodpecker

FROM alpine:latest

RUN apk add --no-cache libc6-compat tzdata

ENV TZ=${TZ:-UTC}

WORKDIR /app

COPY --from=builder /app/woodpecker .

LABEL org.opencontainers.image.source="https://github.com/timberware/woodpecker"
LABEL org.opencontainers.image.description="Woodpecker is a tool that automatically updates DNS A records for specified domain registrars."
LABEL org.opencontainers.image.licenses=GNUGPLv3

CMD ["./woodpecker"]
