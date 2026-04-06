NAME    := cubox-cli
MODULE  := github.com/OLCUBO/cubox-cli
VERSION := $(shell node -p "require('./package.json').version" 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X $(MODULE)/cmd.Version=$(VERSION)

PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64
DIST_DIR  := dist

.PHONY: build install clean build-all release

build:
	go build -ldflags "$(LDFLAGS)" -o $(NAME) .

install:
	go install -ldflags "$(LDFLAGS)" .

clean:
	rm -rf $(NAME) $(DIST_DIR)

build-all: clean
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} ; \
		output=$(NAME) ; \
		if [ "$$GOOS" = "windows" ]; then output=$(NAME).exe; fi ; \
		echo "Building $$GOOS/$$GOARCH..." ; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$$GOOS-$$GOARCH/$$output . ; \
	done

release: build-all
	@mkdir -p $(DIST_DIR)/release
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} ; \
		dir=$(DIST_DIR)/$$GOOS-$$GOARCH ; \
		archive=$(NAME)-$(VERSION)-$$GOOS-$$GOARCH ; \
		if [ "$$GOOS" = "windows" ]; then \
			(cd $$dir && zip ../../$(DIST_DIR)/release/$$archive.zip $(NAME).exe) ; \
		else \
			tar -czf $(DIST_DIR)/release/$$archive.tar.gz -C $$dir $(NAME) ; \
		fi ; \
		echo "Packaged $$archive" ; \
	done
	@echo "Release archives in $(DIST_DIR)/release/"
