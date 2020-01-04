// Huffman 压缩测试
package main

import (
	"dxkite.cn/GoHuffman/huffman"
	"flag"
	"fmt"
	"os"
)

func main() {
	var inputFilename = flag.String("name", "", "the file name be input")
	var help = flag.Bool("help", false, "the file name be input")
	var outputFilename = flag.String("output", "", "the file name be input")
	var mode = flag.String("mode", "encode", "the mode encode/decode")

	flag.Parse()
	if len(os.Args) == 1 || *help {
		flag.Usage()
	}

	switch *mode {
	case "encode", "e":
		if len(*inputFilename) == 0 {
			fmt.Println("please input filename")
			os.Exit(2)
		}
		if len(*outputFilename) == 0 {
			*outputFilename = *inputFilename + ".hdx1"
		}
		err := huffman.EncodeFile(*inputFilename, *outputFilename)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	case "decode", "d":
		if len(*inputFilename) == 0 {
			fmt.Println("please input filename")
			os.Exit(2)
		}
		if len(*outputFilename) == 0 {
			fmt.Println("please input output filename")
			os.Exit(2)
		}
		err := huffman.DecodeFile(*inputFilename, *outputFilename)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("please input mode with encode,e,decode,d")
	}
}
