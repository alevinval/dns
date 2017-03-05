package main

import (
	"flag"
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

	fr := dns.NewReader(f)
	var msg *dns.Msg
	for err == nil {
		msg, err = fr.Read()
		if err == nil {
			debug.PrintMessage(msg)
		}
	}
}
