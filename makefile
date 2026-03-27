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

