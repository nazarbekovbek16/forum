FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && go build -o forum cmd/app/main.go

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app .
EXPOSE 9090
CMD ["/app/forum"]