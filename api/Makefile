OUT := ../dist/appcatalog
PKG := github.com/ovh/lhasa/api
VERSION := $(shell git describe --always --tags --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
PKG_LIST_COMMA := $(shell go list ${PKG}/... | grep -v /vendor/ | paste -s -d, -)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
PACKR = ${GOPATH}/bin/packr
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

$(PACKR):
	go get -u github.com/gobuffalo/packr/packr

install: $(GO_DEP) $(PACKR)
	dep ensure -vendor-only
	packr

server: vet lint install
	packr
	go build -i -o ${OUT} -ldflags="-X main.version=${VERSION}" ${PKG}/cmd/appcatalog

test: install
	@go test ./...

unused:
	codecoroner -ignore vendor funcs ./...

test-coverage:
	@for pkg in ${PKG_LIST}; do go test --coverprofile $(TARGET_DIR)/$${pkg//\//-}.single.cov $${pkg}; done;
	@for pkg in ${PKG_LIST}; do go test --coverprofile $(TARGET_DIR)/$${pkg//\//-}.global.cov --coverpkg=${PKG_LIST_COMMA} $${pkg}; done;
	@$(GO_COV_MERGE) $(TARGET_DIR)/*.global.cov > $(TARGET_DIR)/global.cov
	@go tool cover -html=$(TARGET_DIR)/global.cov -o=$(TARGET_DIR)/cover.html
	@go tool cover -func=$(TARGET_DIR)/global.cov

vet: install
	@go vet ${PKG_LIST}

lint: $(GO_LINT) install
	@for file in ${GO_FILES} ;  do \
		golint -set_exit_status $$file ; \
	done

static: vet lint
	go build -i -v -o ${OUT}-v${VERSION} -tags netgo -ldflags="-extldflags \"-static\" -w -s -X main.version=${VERSION}" ${PKG}/app

run: server
	./${OUT}

live:
	APPCATALOG_CONFIG_FILE=../.config.json APPCATALOG_DEBUG_MODE=true LHASA_WEB_UI_DIR=../dist/webui gin -p 8081 -a 3000 -i --path . --build cmd/appcatalog -- --port=3000

clean:
	@packr clean
	-@rm ${OUT} ${OUT}-v*

.PHONY: run server static vet lint ensure integration-test
