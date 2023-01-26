BIN_PATH=bin
BIN_NAME=perforator
INSTALL_PATH=/usr/local/bin

ifeq (run, $(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif


test:
	go test -v ./...

build:
	@go build -o $(BIN_PATH)/$(BIN_NAME)

.PHONY: run
run: build
	./$(BIN_PATH)/$(BIN_NAME) $(RUN_ARGS)

.PHONY: install
install: build
	cp ./$(BIN_PATH)/$(BIN_NAME) $(INSTALL_PATH)/$(BIN_NAME)

.PHONY: clean
clean:
	rm -rf $(BIN_PATH)/$(BIN_NAME) $(INSTALL_PATH)/$(BIN_NAME)
