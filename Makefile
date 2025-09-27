
APP_NAME=netools
CMD_DIR=./cmd/$(APP_NAME)


run:
	@echo "Running $(APP_NAME)..."
	go run $(CMD_DIR)


.PHONY: run