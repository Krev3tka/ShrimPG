FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /shrimpg ./cmd/passwordManager/main.go


FROM alpine:latest
WORKDIR /root/

COPY --from=builder /shrimpg .
COPY --from=builder /app/migrations ./migrations

CMD ["./shrimpg"]