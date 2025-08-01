BINARY_NAME=ecs_exporter
MAIN_FILE=cmd/main.go

.PHONY: all clear build run stop

clear:
	@echo "Cleaning up binary..."
	@rm -f $(BINARY_NAME)

build:
	@echo "Building binary..."
	@go build -o $(BINARY_NAME) $(MAIN_FILE)

run: build
	@echo "Running in background..."
	@nohup ./$(BINARY_NAME) > /dev/null 2>&1 & echo $$! > $(BINARY_NAME).pid
	@echo "Started $(BINARY_NAME)"

stop:
	@if [ -f $(BINARY_NAME).pid ]; then \
		PID=$$(cat $(BINARY_NAME).pid); \
		echo "Stopping $(BINARY_NAME) (PID $$PID)..."; \
		kill $$PID && rm -f $(BINARY_NAME).pid; \
	else \
		echo "No PID file found. Is $(BINARY_NAME) running?"; \
	fi
