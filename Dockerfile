FROM oven/bun:latest AS css
WORKDIR /app
COPY input.css .
COPY templates/ templates/
COPY internal/ internal/
RUN bunx @tailwindcss/cli -i input.css -o static/output.css

FROM golang:1.25-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=css /app/static/output.css static/output.css
RUN CGO_ENABLED=1 GOOS=linux go build -a -o relay ./cmd/relay

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/relay .
COPY --from=builder /app/templates ./templates
COPY --from=css /app/static ./static
EXPOSE 2323
CMD ["./relay"]
