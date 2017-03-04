package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/go-rfc/dns"
	"github.com/go-rfc/dns/debug"
)

var (
	dumpFile = flag.String("dump", "dump.bin", "dump file path.")
)

func main() {
	flag.Parse()

	f, err := os.OpenFile(*dumpFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Panicf("error opening dump file: %s", err)
	}
	defer f.Close()

	var offset, n int
	err = nil
	data := make([]byte, 512)
	for err != io.EOF {
		if offset >= n-12 {
			offset = 0
			n, err = f.Read(data)
			if err != nil {
				continue
			}
		}
		// TODO: fix reader to check whether message can be read from a given buffer
		// without failing with out of bounds.
		r := dns.NewReader(data[offset:n])
		msg, n2 := r.ReadMessage()
		debug.PrintMessage(msg)
		offset += n2
	}
}
