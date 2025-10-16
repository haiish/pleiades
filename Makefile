APP_NAME = pleiades
BUILD_DIR = build
MAIN_FILE = cmd/main.go
GPRC_PORT = 3552
FIBER_PORT = 8081

.PHONY: run build stop

run:
	go run $(MAIN_FILE)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

stop:
	@echo "🌸 Stopping Fiber on port 8081..."
	@fiber_pid=$$(lsof -i :8081 -t 2>/dev/null); \
	if [ -n "$$fiber_pid" ]; then \
		kill -9 $$fiber_pid; \
		echo "Fiber stopped! (＾▽＾)"; \
	else \
		echo "Fiber not running! (｡•́︿•̀｡)"; \
	fi

	@echo "🌸 Stopping gRPC on port 3552..."
	@grpc_pid=$$(lsof -i :3552 -t 2>/dev/null); \
	if [ -n "$$grpc_pid" ]; then \
		kill -9 $$grpc_pid; \
		echo "gRPC stopped! (≧ω≦)"; \
	else \
		echo "gRPC not running! (｡•́︿•̀｡)"; \
	fi
