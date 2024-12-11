FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o service

FROM alpine:3.20 AS final

COPY --from=builder /app/service /app/service

ENTRYPOINT [ "/app/service" ]