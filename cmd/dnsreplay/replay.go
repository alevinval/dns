package main

import (
	"flag"
	"log"
	"os"
	"io"

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

	fr := dns.NewReader(f)
	var msg *dns.Msg
	for {
		msg, err = fr.Read()
		if err == io.EOF {
			return
		} else if err != nil {
			log.Panicf("error replaying: %s\n", err)
		}
		debug.PrintMessage(msg)
	}
}
