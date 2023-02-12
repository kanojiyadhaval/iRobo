BINARY_NAME=iRobo

all: build

build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)_linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)_darwin_amd64
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)_windows_amd64.exe

clean:
	rm $(BINARY_NAME)_*
