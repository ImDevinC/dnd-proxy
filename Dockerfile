FROM golang:1.23 AS builder

WORKDIR /usr/src/app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

FROM gcr.io/distroless/static

COPY --from=builder /usr/src/app/server /

ENTRYPOINT ["/server"]
