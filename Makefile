BINARY=helloyolo

all:
	@mkdir -p output
	go build -o output/${BINARY} main.go
