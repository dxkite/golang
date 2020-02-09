// Go语言图片隐写
// RGBA通道全部塞入数据
// RGBA会把数据转换成NRGBA照成数据丢失，所以使用NRGBA来写入
// 看源码是个好习惯
package main

import (
	"dxkite.cn/demo/GoPngSteganography/pngio"
	"flag"
	"log"
	"strings"
)

func main() {
	var decode = flag.Bool("decode", false, "decode the input")
	var help = flag.Bool("help", false, "print help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if flag.NArg() > 0 {
		for _, name := range flag.Args() {
			if *decode {
				decodeName, _ := getNameExt(name)
				log.Println("decode", name, "to", decodeName)
				if err := pngio.DecodeFile(name, decodeName); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Println("encode", name, "to", name+".png")
				if err := pngio.EncodeFile(name, name+".png"); err != nil {
					log.Fatal(err)
				}
			}
		}
	} else {
		flag.Usage()
		return
	}
}

func getNameExt(input string) (name, ext string) {
	i := strings.LastIndex(input, ".")
	name, ext = input[:i], input[i+1:]
	return
}
