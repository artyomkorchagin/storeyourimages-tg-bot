FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bot ./cmd/main.go

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/bot /app/bot
EXPOSE 3000
CMD ["/app/bot"]