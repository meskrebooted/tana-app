# --- build stage ---
FROM golang:1.22-alpine AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download 2>/dev/null || true
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o coworking-backend .

# --- run stage ---
FROM alpine:3.20
RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=build /app/coworking-backend .
USER appuser

EXPOSE 8080
ENTRYPOINT ["./coworking-backend"]
