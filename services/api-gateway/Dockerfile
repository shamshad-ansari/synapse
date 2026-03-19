# Stage 1: build
FROM golang:1.23-alpine AS builder

WORKDIR /src/services/api-gateway

COPY services/api-gateway/go.mod services/api-gateway/go.sum ./
RUN go mod download

COPY services/api-gateway/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /app ./cmd/server

# Stage 2: minimal runtime image
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app /app

USER nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]
