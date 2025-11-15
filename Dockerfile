FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /out/app /app/app

COPY migrations ./migrations/
COPY openapi.yml ./

EXPOSE 8080

CMD ["./app"]