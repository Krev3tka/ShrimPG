.PHONY: build run clear

build:
	go build -o /bin/shrimpg .cmd/passwordManager/main.go

run:
	go run .cmd/passwordManager/main.go

clear:
	rm -rf bin/