build: clean get build-bin

# Get application dependencies
get:
	go get ./...

# Compile each cmd app into their own binary
build-bin:
	go build -o bin/$(notdir $(shell pwd)) .
	chmod +x bin/$(notdir $(shell pwd))

# Clean the output files
clean:
	go clean -i ./...
	rm -rf bin

# Run unit and functional tests
test:
	go test -cover -race ./...

# Check construct and style
check: govet golint

# Check for suspicious constructs
govet:
	go tool vet .

# Check for style mistakes
golint:
	golint .
