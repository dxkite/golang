package main

import (
	"dxkite.cn/GoProxy/proxy"
	"log"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	proxy.StartHTTP(":8888")
}
