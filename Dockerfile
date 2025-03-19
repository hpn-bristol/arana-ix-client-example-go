FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ix-go-xapp .

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/ix-go-xapp .
CMD ["./ix-go-xapp"]
