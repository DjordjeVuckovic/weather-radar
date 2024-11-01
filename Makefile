APP_NAME = weather-radar
CMD_DIR = ./cmd/$(APP_NAME)
BUILD_DIR = ./build

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)


.PHONY: fmt
fmt:
	go fmt ./...