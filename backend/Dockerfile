FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod ./
RUN go mod download 2>/dev/null || true
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates poppler-utils
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY --from=builder /app/configs /app/configs
EXPOSE 8080
CMD ["/app/server"]
