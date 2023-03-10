FROM golang:1.19-alpine AS builder

LABEL maintainer="Bek"

WORKDIR /maindir

COPY . .

RUN apk add build-base && go build -o forum cmd/main.go

FROM alpine

WORKDIR /maindir

COPY --from=builder /maindir .

EXPOSE 8000

CMD ["./forum"]