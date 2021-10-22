default: build install

build:
	@echo "Building binary..."
	go build -o bin/statica statica.go

build-all:
	@echo "Building binaries for all platforms..."
	GOOS=windows GOARCH=amd64 go build -o bin/statica-win statica.go
	GOOS=darwin GOARCH=amd64 go build -o bin/statica-osx statica.go
	GOOS=linux GOARCH=amd64 go build -o bin/statica-lin statica.go

install:
	@echo "Installing binary..."
	cp bin/statica ~/bin