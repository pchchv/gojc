# --- Stage 1: Build ---
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o echo-server .

# --- Stage 2: Final lightweight image ---
FROM scratch
WORKDIR /
COPY --from=builder /app/echo-server /echo-server

# COPY your .env file into the root directory so godotenv can find it
COPY .env /.env

EXPOSE 8080
ENTRYPOINT ["/echo-server"]