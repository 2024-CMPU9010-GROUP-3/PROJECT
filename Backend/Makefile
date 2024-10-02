OUTPUT_DIR = ./internal/db

.PHONY: generate clean

all: generate

generate:
	@echo "Generating Go code with sqlc"
	@sqlc generate

clean:
	@rm -rf $(OUTPUT_DIR)/
