BINARY  := urscript-ls
CMD     := ./cmd/urscript-ls

.PHONY: build test vet clean install vscode-ext

build:
	go build -o $(BINARY) $(CMD)

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)

install:
	go install $(CMD)

vscode-ext:
	cd editors/code && npm install && npm run compile
