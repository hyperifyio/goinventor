.PHONY: build clean tidy

GOINVENTOR_SOURCES := ./cmd/goinventor/main.go

all: build

build: goinventor

tidy:
	go mod tidy

goinventor: $(GOINVENTOR_SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o goinventor ./cmd/goinventor

clean:
	rm -f goinventor
