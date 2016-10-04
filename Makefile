VERSION=`git describe --tags 2>/dev/null || echo master`

all: build

build:
	go build -ldflags "-X main.version=${VERSION}"
