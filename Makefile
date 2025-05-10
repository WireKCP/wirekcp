TARGET_EXEC = wirekcp

ifeq ($(OS),Windows_NT)
	FILE_EXTENSION := .exe
else
	FILE_EXTENSION :=
endif

BINARY := main

all: build

build:
	go build -o $(TARGET_EXEC)$(FILE_EXTENSION)

clean:
	rm $(TARGET_EXEC)$(FILE_EXTENSION)

linux:
	GOOS=linux GOARCH=amd64 go build -o $(TARGET_EXEC)

darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(TARGET_EXEC)

darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(TARGET_EXEC)

freebsd:
	GOOS=freebsd GOARCH=amd64 go build -o $(TARGET_EXEC)

windows:
	GOOS=windows GOARCH=amd64 go build -o $(TARGET_EXEC)$(FILE_EXTENSION)

.PHONY: all build clean