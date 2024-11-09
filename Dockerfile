FROM golang:1.21 AS builder

WORKDIR /src
COPY . .
COPY .env .env
RUN apt-get update && apt-get install -y ca-certificates
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"' .

FROM scratch
COPY --from=builder /app /app
COPY --from=builder /src/.env .env

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


EXPOSE 8080

ENTRYPOINT ["/app"]
