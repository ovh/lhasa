ARCH = $(shell uname | tr '[:upper:]' '[:lower:]')
TARGET_DIR = ./dist
COMPOSE_BIN = $(TARGET_DIR)/docker-compose
API_BIN = $(TARGET_DIR)/appcatalog
VENOM_BIN = $(TARGET_DIR)/venom.$(ARCH)
UI_BASE_HREF = /
export UI_BASE_HREF

all: api webui

$(TARGET_DIR):
	$(info Creating $(TARGET_DIR) directory)
	@mkdir -p $(TARGET_DIR)

$(VENOM_BIN): $(TARGET_DIR)
	$(info Installing venom... for $(ARCH))
	curl -L -o $(VENOM_BIN) https://github.com/ovh/venom/releases/download/v0.17.0/venom.$(ARCH)-amd64
	@chmod +x $(VENOM_BIN)

$(COMPOSE_BIN): $(TARGET_DIR)
	$(info Installing docker-compose...)
	@curl -L https://github.com/docker/compose/releases/download/1.17.0/docker-compose-`uname -s`-`uname -m` -o $(COMPOSE_BIN)
	@chmod +x $(COMPOSE_BIN)

$(API_BIN):
	$(MAKE) -C api server

api:
	$(MAKE) -C api
	cp scripts/appcatalog.sh dist
	cp samples/mycompany.sh dist

webui:
	$(MAKE) -C webui

test:
	$(MAKE) -C api test
	$(MAKE) -C webui test

run: all
	cd dist && ./appcatalog --config=../.config.json --auto-migrate

clean:
	$(MAKE) -C api clean
	$(MAKE) -C webui clean

integration-test: $(COMPOSE_BIN) $(VENOM_BIN) $(API_BIN)
	$(COMPOSE_BIN) up -d
	sleep 10;
	{ ./${API_BIN} ${DEBUG} --auto-migrate & }; \
	pid=$$!; \
	sleep 5; \
	APP_HOST=http://localhost:8081 $(VENOM_BIN) run --strict --output-dir=$(TARGET_DIR) tests/; \
	r=$$?; \
	kill $$pid; \
	./${API_BIN} ${DEBUG} migrate down; \
	$(COMPOSE_BIN) down; \
	exit $$r

sample-test: $(VENOM_BIN) $(API_BIN)
	APP_HOST=http://localhost:8081 $(VENOM_BIN) run --log debug --format xml --output-dir=$(TARGET_DIR) tests/10-samples-v1.yml && cat $(TARGET_DIR)/test_results.xml

.PHONY: all test run clean integration-test api webui
