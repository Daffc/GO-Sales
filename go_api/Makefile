.PHONY: all clean swag build

# Project folders
SWAGGER_DIR = ./docs
CMD_DIR = ./cmd
HANDLER_DIR = ./api/handler/
DTO_DIR = ./domain/dto

# Output program name
BINARY_NAME = server

all: swag build

# Generate API documentation with
swag:
	swag init --output docs --dir $(CMD_DIR),$(HANDLER_DIR),$(DTO_DIR)

# Compiles project
build:
	go build -o $(BINARY_NAME) $(CMD_DIR)/main.go

clean:
	rm -rf $(SWAGGER_DIR)
	rm -f $(BINARY_NAME)