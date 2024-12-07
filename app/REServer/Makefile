PROG=bin/REServer


SRCS=./

# installation manual
INSTALL_PREFIX=/usr/local

# Configure the installation directory
CONF_INSTALL_PREFIX=/usr/local/REServer

# git commit hash
COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")

# compile date
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')

# compile condition
CFLAGS = -ldflags "-s -w -X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\""

all:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	go build $(CFLAGS) -o $(PROG) $(SRCS)

linux_x86_64:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	CGO_ENABLED=0  GOOS=linux  GOARCH=amd64 go build $(CFLAGS) -o $(PROG)_linux_x86_64 $(SRCS)

win64:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	CGO_ENABLED=0 GOOS=windows  GOARCH=amd64 go build $(CFLAGS) -o $(PROG)_win64.exe $(SRCS)

macos_x86_64:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64  go build $(CFLAGS) -o $(PROG)_macos_x86_64 $(SRCS)

# Compile the race version
race:
	if [ ! -d "./bin/" ]; then \
    	mkdir bin; \
    	fi
	go build $(CFLAGS) -race -o $(PROG) $(SRCS)

# release version
RELEASE_DATE = $(shell date '+%Y%m%d%H%M%S')
RELEASE_VERSION = $(shell git rev-parse --short HEAD || echo "GitNotFound")
RELEASE_DIR=release_bin
RELEASE_BIN_NAME=REServer
release:
	if [ ! -d "./$(RELEASE_DIR)/$(RELEASE_DATE)_$(RELEASE_VERSION)" ]; then \
	mkdir -p ./$(RELEASE_DIR)/$(RELEASE_DATE)_$(RELEASE_VERSION); \
	fi
	go build $(CFLAGS) -o $(RELEASE_DIR)/$(RELEASE_DATE)_$(RELEASE_VERSION)/$(RELEASE_BIN_NAME) $(SRCS)


install:
	cp $(PROG) $(INSTALL_PREFIX)/bin
	
	if [ ! -d "${CONF_INSTALL_PREFIX}" ]; then \
	mkdir $(CONF_INSTALL_PREFIX); \
	fi
	
	cp -R etc/* $(CONF_INSTALL_PREFIX)

clean:
	rm -rf ./bin
	
	rm -rf $(INSTALL_PREFIX)/bin/REServer
	
	rm -rf $(CONF_INSTALL_PREFIX)
