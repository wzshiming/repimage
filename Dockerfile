FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o repimage .

FROM alpine
COPY --from=builder /app/repimage /repimage
COPY ./certs /certs
