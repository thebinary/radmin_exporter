SRCS = $(wildcard *.go) $(wildcard exporters/*.go) $(wildcard libradmin/*.go)

.PHONY: build
build: build/radmin_exporter
build/radmin_exporter: $(SRCS)
	go build -o $@
	upx $@

all: build/linux/amd64 build/darwin/amd64

.PHONY: build/linux/amd64
build/linux/amd64: build/radmin_exporter-linux-amd64
build/radmin_exporter-linux-amd64: $(SRCS)
	GOOS=linux GOARCH=amd64 go build -o $@ *.go
	upx $@

.PHONY: build/darwin/amd64
build/darwin/amd64: build/radmin_exporter-darwin-amd64
build/radmin_exporter-darwin-amd64: $(SRCS)
	GOOS=darwin GOARCH=amd64 go build -o $@ *.go
	upx $@

.PHONY: clean
clean:
	rm -rf build
