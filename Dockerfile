FROM golang:alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download || true

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/api/main.go || exit 0

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api* ./api

EXPOSE 8080

CMD ["./api"]
