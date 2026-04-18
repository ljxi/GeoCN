FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server .

FROM alpine:3.21

WORKDIR /app
COPY --from=builder /server .
COPY division_code/ division_code/
COPY db/ db/

EXPOSE 80
ENTRYPOINT ["/app/server"]
