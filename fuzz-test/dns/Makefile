WORKDIR = .fuzz/
BIN = ./dns-fuzz.zip

build:
	go-fuzz-build github.com/go-rfc/dns/fuzz-test/dns

test:
	go-fuzz -bin=$(BIN) -workdir=$(WORKDIR)

clean:
	rm -rf $(WORKDIR)
	rm -f $(BIN)

all: clean build test
