FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o findopedia ./cmd/findopedia/

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /build/findopedia .
COPY migrations/ migrations/
EXPOSE 8080
CMD ["./findopedia", "-migrationsDir", "migrations"]
