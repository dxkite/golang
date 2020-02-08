package main

import (
	"dxkite.cn/demo/GoVideo2Image/upload"
	"dxkite.cn/demo/GoVideo2Image/video"
	"flag"
	"log"
	"os"
	"path"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	var input = flag.String("input", "", "the mp4 file")
	var t = flag.String("type", "ali", "the type to upload image")
	var time = flag.Int("time", 30, "time per segment")
	var outputDir = flag.String("temp", "tmp", "the template dir")
	var outputIndex = flag.String("output", "output.m3u8", "the output m3u8 index")
	var binary = flag.String("bin", "ffmpeg", "ffmpeg binary command")
	var ext = flag.String("ext", "jpg", "image extension")

	var help = flag.Bool("help", false, "the file name be input")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	tempIndex := path.Join(*outputDir, "output.m3u8")
	cvt := video.NewSimpleConverter(*binary, *ext, *outputDir, os.Stdout)
	if er := cvt.Convert("go-t-video-", *input, tempIndex, *time); er != nil {
		log.Fatal(er)
	}

	if er := video.MakeM3u8("go-t-video-", tempIndex, *outputIndex, *outputDir, func(name string, data []byte) (url string, err error) {
		log.Println("upload", name, "...")
		re, err := upload.Upload(*t, &upload.FileObject{
			Name: name,
			Data: data,
		})
		if err != nil {
			log.Println("upload error", name)
			return "", err
		}
		return re.Url, nil
	}); er != nil {
		log.Fatal(er)
	}
	log.Println("success")
}
