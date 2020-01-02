package main

import (
	"dxkite.cn/GoProxy/web"
	"flag"
	"log"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	//proxy.StartHTTP(":8888")
	var listen = flag.String("listen", ":8080", "the listen port")
	var pacFile = flag.String("pac", "pac.txt", "the pac.txt")
	var proxy = flag.String("proxy", "127.0.0.1:1080", "the proxy")
	var help = flag.Bool("help", false, "the file name be input")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	web.StartPacServe(*listen, *pacFile, *proxy)
}
