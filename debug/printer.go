package debug

import (
	"encoding/json"
	"fmt"

	"github.com/go-rfc/dns/parse"
)

func PrintMessage(msg *parse.Msg) {
	d, _ := json.Marshal(msg.Header)
	fmt.Printf("HEADER:\n  %s\n", d)
	for i, q := range msg.Queries {
		d, _ = json.Marshal(q)
		fmt.Printf("QUERY %d:\n  %s\n", i, d)
	}
}
