FROM golang:1.25-alpine
WORKDIR /app
COPY . .
RUN ls -R /
RUN go build -o shrimpg ./cmd/passwordManager/main.go
CMD ["./shrimpg"]