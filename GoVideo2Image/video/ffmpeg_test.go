package video

import (
	"log"
	"os"
	"testing"
)

func TestFFmpegConvert_Convert(t *testing.T) {
	cvt := NewSimpleConverter(`E:\software\ffmpeg-20200124-e931119-win64-static\bin\ffmpeg.exe`, "jpg", "output", os.Stdout)
	if er := cvt.Convert("cdn-video-", `E:\01 诞生!救世主猪猪侠_高清.mp4`, `output\test.m3u8`, 30); er != nil {
		log.Print(er)
	}
}
