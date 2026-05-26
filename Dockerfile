# =========================================
# Stage 1: Build
# =========================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o otp-service ./cmd/app

# =========================================
# Stage 2: Runtime
# =========================================
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/otp-service .

# Copy config folder
COPY --from=builder app/config ./config

EXPOSE 9999

CMD ["./otp-service"]