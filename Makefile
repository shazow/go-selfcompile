BUILD = ./build
VERSION = $(shell git describe --long --tags --dirty --always)
SOURCES = $(wildcard *.go **/*.go)

all: $(BUILD)/go-selfcompile

example-abinary: $(BUILD)/example-abinary

example-aplugin: example-abinary
	$(BUILD)/example-abinary --plugin "github.com/shazow/go-selfcompile/example/aplugin"

test:
	go test .

$(BUILD)/go-selfcompile: $(SOURCES)
	go build -ldflags "-X main.version='$(VERSION)'" -o "$@" ./go-selfcompile/...

example/abinary/bindata_selfcompile.go: $(BUILD)/go-selfcompile
	cd example/abinary/ && ../../$< --debug --skip-source

$(BUILD)/example-abinary: example/abinary/bindata_selfcompile.go
	go build -o "$@" ./example/abinary/...

$(BUILD):
	mkdir $@

clean:
	rm -r $(BUILD)
