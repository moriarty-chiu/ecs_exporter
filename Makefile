BINARY_NAME=ecs_exporter
MAIN_FILE=cmd/main.go
LOG_FILE=run.log

.PHONY: all clear build run

clear:
	@echo "Cleaning up binary..."
	@rm -f $(BINARY_NAME)

build:
	@echo "Building binary..."
	@go build -ldflags="-s -w" -a -o $(BINARY_NAME) $(MAIN_FILE)

run: build
	@echo "Running in background..."
	@nohup ./$(BINARY_NAME) > $(LOG_FILE) 2>&1 &
	@echo "Started $(BINARY_NAME), logging to $(LOG_FILE)"
