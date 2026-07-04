FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o vigilante .

FROM alpine:3.23
RUN apk add --no-cache bash
WORKDIR /app
COPY --from=builder /app/vigilante .
COPY config.yaml scripts/ ./scripts/
CMD ["./vigilante"]
