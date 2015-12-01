package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/levenlabs/go-srvclient"
	"github.com/levenlabs/golib/rpcutil"
	"net/url"
)

func main() {
	bin := os.Args[0]
	argv := os.Args[1:]
	if len(argv) != 3 {
		fmt.Printf("Usage: %s <url> <service.method> <parameters>\n\n", bin)
		fmt.Printf(`Example: %s 'http://localhost:5555/api' 'Foo.DoAThing' '{"foo":"bar"}'`, bin)
		fmt.Print("\n\n")
		return
	}

	u, method, body := argv[0], argv[1], argv[2]
	if uo, err := url.Parse(u); err == nil {
		uo.Host = srvclient.MaybeSRV(uo.Host)
		u = uo.String()
	}
	var ret interface{}
	err := rpcutil.JSONRPC2RawCall(u, &ret, method, body)
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := json.MarshalIndent(ret, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(out))
}
