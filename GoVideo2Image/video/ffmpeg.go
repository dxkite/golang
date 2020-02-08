package video

import (
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
)

type FFmpegConvert struct {
	Binary     string
	OutputPath string
	Stdout     io.Writer
	Stderr     io.Writer
	Extension  string
}

func NewSimpleConverter(binary, ext, output string, writer io.Writer) *FFmpegConvert {
	return &FFmpegConvert{
		Binary:     binary,
		OutputPath: output,
		Stdout:     writer,
		Stderr:     writer,
		Extension:  ext,
	}
}

func (cvt *FFmpegConvert) Convert(prefix, input, output string, segmentTime int) error {
	_ = os.MkdirAll(cvt.OutputPath, os.ModePerm)
	cmd := exec.Command(cvt.Binary, "-i", input,
		"-f", "segment",
		"-segment_time", strconv.Itoa(segmentTime),
		"-segment_format", "mpegts",
		"-segment_list", output,
		"-c", "copy",
		"-bsf:v", "h264_mp4toannexb", "-map", "0", path.Join(cvt.OutputPath, prefix+"%04d."+cvt.Extension))
	cmd.Stdout = cvt.Stdout
	cmd.Stderr = cvt.Stderr
	er := make(chan error)
	go func() {
		er <- cmd.Run()
	}()
	return <-er
}
