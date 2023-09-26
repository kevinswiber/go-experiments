FROM golang:bookworm as builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY goserver/main.go .
RUN go build -o goserver ./main.go

FROM debian:12-slim

WORKDIR /app
COPY --from=builder /build/goserver .

EXPOSE 8080
ENV HOST_ADDR="0.0.0.0"
CMD ["/app/goserver"]
