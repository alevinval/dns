package dns

import (
	"bytes"
	"fmt"

	"github.com/go-rfc/dns"
	"github.com/go-rfc/dns/debug"
)

func Fuzz(input []byte) int {
	msg, n, err := dns.UnpackMsg(input, 0)
	if err != nil {
		if msg != nil {
			panic("msg should be nil on error")
		}
		return 0
	}

	output := dns.PackMsg(msg)
	input = input[:n]
	if !bytes.Equal(input, output) {
		fmt.Printf("%b\n", input)
		fmt.Printf("%b\n", output)
		reUnPackedMsg, _, _ := dns.UnpackMsg(output, 0)
		debug.PrintMessage(msg)
		debug.PrintMessage(reUnPackedMsg)
		panic("input is not equal to the output")
	}

	return 1
}
