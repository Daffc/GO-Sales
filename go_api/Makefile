.PHONY: all clean swag

# Project folders
SWAGGER_DIR = ./docs
CMD_DIR = ./cmd
HANDLER_DIR = ./internal/handler/
USECAS_DIR = ./internal/domain/usecase/

# Output program name
BINARY_NAME = server

all: swag build

# Generate API documentation with
swag:
	swag init --output docs --dir $(CMD_DIR),$(HANDLER_DIR),$(USECAS_DIR)

# Compiles project
build:
	go build -o $(BINARY_NAME) $(CMD_DIR)/main.go

# Migrate Database
migrate:
	cd scripts; ./goose.sh up;

clean:
	rm -rf $(SWAGGER_DIR)
	rm -f $(BINARY_NAME)