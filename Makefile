default: build install

build:
	@echo "Building binary..."
	go build -o bin/statika statika.go

build-all:
	@echo "Building binaries for all platforms..."
	GOOS=windows GOARCH=amd64 go build -o bin/statika-win statika.go
	GOOS=darwin GOARCH=amd64 go build -o bin/statika-osx statika.go
	GOOS=linux GOARCH=amd64 go build -o bin/statika-lin statika.go

install:
	@echo "Installing binary..."
	cp bin/statika ~/bin