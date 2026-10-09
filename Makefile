BINARY  := urscript-ls
CMD     := ./cmd/urscript-ls
VERSION ?= dev

TARGETS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 windows-arm64
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build test vet clean install vscode-ext release-local

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) $(CMD)

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
	rm -rf dist/
	rm -f editors/code/*.vsix
	rm -rf editors/code/server/

install:
	go install -ldflags="$(LDFLAGS)" $(CMD)

vscode-ext:
	cd editors/code && npm install && npm run compile

cross-compile:
	@mkdir -p dist
	@for target in $(TARGETS); do \
		os=$$(echo $$target | cut -d- -f1); \
		arch=$$(echo $$target | cut -d- -f2); \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		echo "Building $$os/$$arch..."; \
		GOOS=$$os GOARCH=$$arch go build -ldflags="$(LDFLAGS)" \
			-o "dist/$(BINARY)-$$os-$$arch$$ext" $(CMD); \
	done

release-local: cross-compile vscode-ext
	@cd editors/code && \
	for pair in linux-x64:linux-amd64 linux-arm64:linux-arm64 \
		darwin-x64:darwin-amd64 darwin-arm64:darwin-arm64 \
		win32-x64:windows-amd64 win32-arm64:windows-arm64; do \
		vscode_target=$$(echo $$pair | cut -d: -f1); \
		go_target=$$(echo $$pair | cut -d: -f2); \
		os=$$(echo $$go_target | cut -d- -f1); \
		mkdir -p server; \
		if [ "$$os" = "windows" ]; then \
			cp "../../dist/$(BINARY)-$${go_target}.exe" server/$(BINARY).exe; \
		else \
			cp "../../dist/$(BINARY)-$${go_target}" server/$(BINARY); \
			chmod +x server/$(BINARY); \
		fi; \
		npx vsce package --target "$$vscode_target" -o "$(BINARY)-$$vscode_target.vsix"; \
		rm -rf server; \
	done
	@echo "Built .vsix files:"
	@ls -la editors/code/*.vsix
