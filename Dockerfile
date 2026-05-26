# =========================================
# Stage 1: Build
# =========================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Giảm parallel build để đỡ ăn CPU
ENV GOMAXPROCS=2

RUN CGO_ENABLED=0 \
    go build \
    -p 2 \
    -trimpath \
    -ldflags="-s -w" \
    -o otp-service \
    ./cmd/app

# =========================================
# Stage 2: Runtime
# =========================================
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/otp-service .
COPY --from=builder /app/config ./config

EXPOSE 9999

CMD ["./otp-service"]