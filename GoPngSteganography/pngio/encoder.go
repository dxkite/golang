package pngio

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
)

// 解码
func Decode(r io.Reader, w io.Writer) error {
	i, err := png.Decode(r)
	if err != nil {
		return err
	}
	switch i := i.(type) {
	case *image.NRGBA:
		p := &ImagePack{Image: i}
		if _, err := io.Copy(w, p); err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("unknown color mode： image.RGBA"))
	}
	return nil
}

// 编码
func Encode(w io.Writer, r io.Reader) error {
	p := NewPack(1024)
	if _, err := io.Copy(p, r); err != nil {
		return err
	}
	return png.Encode(w, p.Resize())
}

func EncodeFile(input, output string) error {
	fi, ei := os.OpenFile(input, os.O_RDONLY, os.ModePerm)
	if ei != nil {
		return ei
	}
	fo, eo := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if eo != nil {
		return eo
	}
	return Encode(fo, fi)
}

func DecodeFile(input, output string) error {
	fi, ei := os.OpenFile(input, os.O_RDONLY, os.ModePerm)
	if ei != nil {
		return ei
	}
	fo, eo := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if eo != nil {
		return eo
	}
	return Decode(fi, fo)
}
