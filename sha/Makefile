# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GO_PROXY=https://goproxy.io

all: clean linux macos win

clean:
	rm -f ./sha1-*

win:
	GOOS=windows GOARCH=amd64 $(GOBUILD)  -o ./sha1-win

macos:
	GOOS=darwin GOARCH=amd64 $(GOBUILD)  -o ./sha1-mac

linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD)  -o ./sha1-linux