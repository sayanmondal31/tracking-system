FROM golang:1.26-alpine AS builder

WORKDIR /app 

# Copy dependencies
COPY go.mod go.sum ./
RUN  go mod download

COPY . .

# Compile a static binary optimized for production (no debug info, static links)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api-gateway main.go

FROM golang:1.26-alpine AS runner
WORKDIR /app

COPY --from=builder /app/api-gateway .

EXPOSE 3000

CMD ["./api-gateway"]
