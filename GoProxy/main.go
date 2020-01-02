package main

import (
	"flag"
	"log"
	"os"
	"time"

	"dxkite.cn/GoProxy/proxy"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	var addr = flag.String("serve", "", "the remote address")
	var listen = flag.String("listen", "", "the listen address")
	var http_proxy = flag.String("http_proxy", "", "the second http_proxy")

	var pemFile = flag.String("pem_file", "conf/server.pem", "the certFile")
	var keyFile = flag.String("key_file", "conf/server.key", "the keyFile")

	var help = flag.Bool("help", false, "the file name be input")
	var xor = flag.Int("key", 22, "the xor key")

	flag.Parse()
	if len(os.Args) == 1 || *help {
		flag.Usage()
		return
	}

	if len(*addr) > 0 {
		log.Println("client mode")
		// 代理链接到加密代理
		proxy.StartHTTPWrapperConnect(*listen, nil,
			proxy.NewTLSConnect(*addr, 3*time.Second).SetWrapper(proxy.NewXORWrapper(byte(*xor))))

	} else {
		log.Println("server mode")
		var connector proxy.Connector
		if len(*http_proxy) > 0 {
			connector = proxy.NewHTTPConnect(*http_proxy, 3*time.Second)
		}
		proxy.StartHTTPWrapperConnectListen(*listen, proxy.NewXORWrapper(byte(*xor)),
			connector,
			proxy.NewTLSListen(*pemFile, *keyFile))
	}
}
