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
	var input = flag.String("input", "", "the input")
	var output = flag.String("output", "", "the output")

	var help = flag.Bool("help", false, "print help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if len(*input) > 0 {
		if len(*output) == 0 {
			*output = *input + ".png"
		}
		if *decode {
			log.Println("decode", *input, "to", *output)
			if err := pngio.DecodeFile(*input, *output); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("encode", *input, "to", *output)
			if err := pngio.EncodeFile(*input, *output); err != nil {
				log.Fatal(err)
			}
		}
	} else if flag.NArg() > 0 {
		for _, name := range flag.Args() {
			filename, ext := getNameExt(name)
			if ext == "png" {
				log.Println("decode", name, "to", filename)
				if err := pngio.DecodeFile(name, filename); err != nil {
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
