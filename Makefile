SRCS = $(wildcard *.go) $(wildcard exporters/*.go) $(wildcard libradmin/*.go)

all: linux darwin

.PHONY: linux
linux: build/radmin_exporter-linux-amd64
build/radmin_exporter-linux-amd64: $(SRCS)
	GOOS=linux GOARCH=amd64 go build -o $@ *.go
	upx $@

.PHONY: darwin
darwin: build/radmin_exporter-darwin-amd64
build/radmin_exporter-darwin-amd64: $(SRCS)
	GOOS=darwin GOARCH=amd64 go build -o $@ *.go
	upx $@

.PHONY: clean
clean:
	rm -rf build
