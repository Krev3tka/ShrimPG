BINARY_NAME=shrimpg
BUILD_DIR=bin

.PHONY: build clear run

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o bin/shrimpg cmd/passwordManager/main.go

clear:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

run: build
	@echo "Running..."
	./$(BUILD_DIR)/$(BINARY_NAME)


deploy:
	@echo "Cross-compiling for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/shrimpg-linux cmd/passwordManager/main.go
	@echo "Uploading to server..."
	scp bin/shrimpg-linux $(SERVER_USER)@$(SERVER_IP):/root/shrimpg
	@echo "Done! Connect via SSH and run ./shrimpg"

