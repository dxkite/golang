package video

import (
	"os"
	"testing"
)

func TestMakeM3u8(t *testing.T) {
	outputIndex := `output\test.m3u8`
	uploadIndex := `output\upload.m3u8`
	outputDir := "output"
	outPrefix := "cdn-video"
	outTimeRange := 30

	cvt := NewSimpleConverter(`E:\software\ffmpeg-20200124-e931119-win64-static\bin\ffmpeg.exe`, "jpg", outputDir, os.Stdout)
	if er := cvt.Convert(outPrefix, `E:\01 诞生!救世主猪猪侠_高清.mp4`, outputIndex, outTimeRange); er != nil {
		t.Error(er)
		return
	}

	if er := MakeM3u8(outPrefix, outputIndex, uploadIndex, outputDir, func(name string, data []byte) (url string, err error) {
		return "upload-" + name, nil
	}); er != nil {
		t.Error(er)
	}
}
