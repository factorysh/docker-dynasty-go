GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

binary: bin
	go build -ldflags "-X gitlhub.com/factorysh/docker-dynasty-go/version.version=$(GIT_VERSION)" \
		-o bin/dynasty

bin:
	mkdir -p bin

test:
	go test -v gitlhub.com/factorysh/docker-dynasty-go/dynasty