package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/levenlabs/go-llog"
	"github.com/levenlabs/go-srvclient"
	"github.com/levenlabs/golib/rpcutil"
	"github.com/mediocregopher/lever"
)

func main() {
	l := lever.New("rpc-cli", &lever.Opts{
		HelpHeader:         "Usage: rpc-cli [options] <url> <service.method> <parameters>\n",
		DisallowConfigFile: true,
	})
	l.Add(lever.Param{
		Name:        "--llog",
		Description: "Instead of outputting the return as json, output a llog-style message with the returned json object as the key/value params",
		Flag:        true,
	})
	l.Parse()

	argv := l.ParamRest()
	if len(argv) != 3 {
		fmt.Print(l.Help())
		exit(0)
	}

	u, method, body := argv[0], argv[1], argv[2]
	if uo, err := url.Parse(u); err == nil {
		uo.Host = srvclient.MaybeSRV(uo.Host)
		u = uo.String()
	}
	var ret interface{}
	err := rpcutil.JSONRPC2RawCall(u, &ret, method, body)
	if err != nil {
		llog.Error("error calling rpc method", llog.KV{"err": err})
		exit(1)
	}

	if l.ParamFlag("--llog") {
		retm, ok := ret.(map[string]interface{})
		if !ok {
			llog.Error("return value not a json object, can't llog")
			exit(1)
		}
		llog.Info("output from "+method, llog.KV(retm))
		exit(0)
	}

	out, err := json.MarshalIndent(ret, "", "    ")
	if err != nil {
		llog.Error("error marshalling json", llog.KV{"err": err})
		exit(1)
	}

	fmt.Println(string(out))
}

func exit(i int) {
	time.Sleep(100 * time.Millisecond)
	os.Exit(i)
}
