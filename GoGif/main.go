package main

import (
	"flag"
	"fmt"
	"image"
	"dxkite.cn/demo/GoGif/gogif"
	"log"
	"strings"
	_ "image/png"
	_ "image/jpeg"
)

func getNameExt(input string) (name, ext string ){
	i := strings.LastIndex(input, ".")
	name, ext = input[:i], input[i+1:]
	return
}

func getOutputName(i,max int, input, output string)  string {
	var name, ext string
	ext = "gif"
	if len(output) > 0 {
		name, _ = getNameExt(output)
		if max > 1 {
			name = fmt.Sprintf("%s-%d", name, i)
		}
	} else {
		name, _ = getNameExt(input)
	}
	return name+"."+ext
}

func main() {
	var output = flag.String("output", "", "the output file")
	var help = flag.Bool("help", false, "print help")
	var width = flag.Int("width", 32, "the image width")
	var delay = flag.Int("delay", 8, "the gif delay 100th of second")
	var height = flag.Int("height", 32, "the image height")
	flag.Parse()
	if *help || flag.NArg() < 1 {
		flag.Usage()
		return
	}
	for i, input := range flag.Args() {
		outputName := getOutputName(i, flag.NArg(), input, *output)
		err := gogif.MakeGif(input, outputName, image.Rect(0, 0, *width, *height), *delay)
		if err != nil {
			log.Println("make gif error", input, err)
		} else {
			log.Println("make gif", input, outputName)
		}
	}
}
