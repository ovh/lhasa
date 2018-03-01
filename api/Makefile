OUT := ../dist/appcatalog
PKG := github.com/ovh/lhasa/api
VERSION := $(shell git describe --always --tags --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
GO_LINT = ${GOPATH}/bin/golint
GO_COV_MERGE = ${GOPATH}/bin/gocovmerge
GO_GOVERALLS = ${GOPATH}/bin/goveralls
GO_GO2XUNIT = ${GOPATH}/bin/go2xunit
GO_DEP = ${GOPATH}/bin/dep
TARGET_DIR = ../dist

all: server

$(TARGET_DIR):
	$(info create $(TARGET_DIR) directory)
	@mkdir -p $(TARGET_DIR)

$(GO_LINT):
	go get -u github.com/golang/lint/golint

$(GO_COV_MERGE):
	go get -u github.com/wadey/gocovmerge

$(GO_GO2XUNIT):
	go get -u github.com/tebeka/go2xunit

$(GO_DEP):
	go get -u github.com/golang/dep/cmd/dep

install: $(GO_DEP)
	dep ensure

server: vet lint install
	go build -i -o ${OUT} -ldflags="-X main.version=${VERSION}" ${PKG}/cmd/appcatalog

test: install
	@go test -short ${PKG_LIST}

test-coverage: $(GO_COV_MERGE) $(TARGET_DIR)
	@$(GO_COV_MERGE) ./**/*.coverprofile > $(TARGET_DIR)/cover.out
	@go tool cover -html=$(TARGET_DIR)/cover.out -o=$(TARGET_DIR)/cover.html

integration-test: server
	{ ./${OUT} --auto-migrate & }; \
	pid=$$!; \
	sleep 3; \
	APP_HOST=http://localhost:8081 venom run --output-dir=$(TARGET_DIR) tests/; \
	r=$$?; \
	kill $$pid; \
	exit $$r

vet:
	@go vet ${PKG_LIST}

lint: $(GO_LINT)
	@for file in ${GO_FILES} ;  do \
		golint -set_exit_status $$file ; \
	done

static: vet lint
	go build -i -v -o ${OUT}-v${VERSION} -tags netgo -ldflags="-extldflags \"-static\" -w -s -X main.version=${VERSION}" ${PKG}/app

run: server
	./${OUT}

clean:
	-@rm ${OUT} ${OUT}-v*

.PHONY: run server static vet lint ensure integration-test
