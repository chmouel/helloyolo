BINARY=helloyolo

all:
	@mkdir -p output
	go build -o output/${BINARY} *.go
	@chmod +x output/${BINARY}

install:
	go install .


clean:
	- rm -f output/${BINARY}
