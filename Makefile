BINARY_NAME=iRobo
TAG = $(shell git describe --abbrev=0 --tags)
NEXT_TAG = v$(shell echo $(TAG) | awk -F. '{$$NF+=1; OFS="."; print $$0}')

all: build

build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)_linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)_darwin_amd64
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)_windows_amd64.exe

release:
	@echo "Tagging release as $(NEXT_TAG)"
	git tag $(NEXT_TAG)
	git push origin $(NEXT_TAG)

clean:
	rm $(BINARY_NAME)_*
