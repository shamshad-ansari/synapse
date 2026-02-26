# Build stage
FROM golang:1.23-alpine AS build

WORKDIR /app

# Needed for fetching modules in some environments
RUN apk add --no-cache git ca-certificates

COPY services/api-gateway/go.mod services/api-gateway/go.sum ./
RUN go mod download

COPY services/api-gateway/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api-gateway ./cmd/server

# Run stage
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /bin/api-gateway /api-gateway

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/api-gateway"]