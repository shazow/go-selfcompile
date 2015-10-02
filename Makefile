BUILD = ./build
VERSION = $(shell git describe --long --tags --dirty --always)
SOURCES = $(wildcard *.go **/*.go)

$(BUILD)/go-selfcompile: $(SOURCES)
	go build -ldflags "-X main.version='$(VERSION)'" -o "$@" ./go-selfcompile/...

clean:
	rm -r $(BUILD)

$(BUILD):
	mkdir $@
