BUILD = $$PWD/build
VERSION = $(shell git describe --long --tags --dirty --always)
SOURCES = $(wildcard *.go **/*.go)

all: $(BUILD)/go-selfcompile

example-abinary: $(BUILD)/example-abinary

example-aplugin: example-abinary
	$(BUILD)/example-abinary --plugin "github.com/shazow/go-selfcompile/example/aplugin"

test:
	go test .

$(BUILD)/go-selfcompile: $(BUILD) $(SOURCES)
	go build -ldflags "-X main.version='$(VERSION)'" -o "$@" ./go-selfcompile/...

example/abinary/bindata_selfcompile.go: $(BUILD)/go-selfcompile
	PATH=$(BUILD):$$PATH go generate ./example/abinary/

$(BUILD)/example-abinary: example/abinary/bindata_selfcompile.go
	go build -o "$@" ./example/abinary/...

$(BUILD):
	mkdir -p $@

clean:
	rm -fr $(BUILD)
	rm -f example/abinary/{bindata_selfcompile.go,abinary}
