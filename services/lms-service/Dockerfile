# Stage 1: build
FROM golang:1.23-alpine AS builder

WORKDIR /src/services/lms-service

COPY services/lms-service/go.mod services/lms-service/go.sum ./
RUN go mod download

COPY services/lms-service/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /app ./cmd/server

# Stage 2: minimal runtime image
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app /app

USER nonroot
EXPOSE 8081
ENTRYPOINT ["/app"]
