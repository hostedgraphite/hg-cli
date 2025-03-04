PROJECT_NAME := hg-cli

GO_BUILD := go build -o

# Build targets
build: build-linux build-windows build-mac

build-linux:
	GOOS=linux GOARCH=amd64 $(GO_BUILD) $(OUTPUT_DIR)/$(PROJECT_NAME)-linux-amd64

build-windows:
	GOOS=windows GOARCH=amd64 $(GO_BUILD) $(OUTPUT_DIR)/$(PROJECT_NAME)-windows-amd64.exe

build-mac:
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) $(OUTPUT_DIR)/$(PROJECT_NAME)-darwin-amd64

clean:
	rm -rf $(OUTPUT_DIR)

release-snapshot:
	goreleaser release --snapshot --clean

.PHONY: build build-linux build-windows build-mac clean release-snapshot
