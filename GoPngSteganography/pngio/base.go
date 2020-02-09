package pngio

import (
	"image"
	"io"
	"math"
)

type ImagePack struct {
	Image  *image.NRGBA
	offset int
}

func (p *ImagePack) Write(b []byte) (n int, err error) {
	n = p.offset
	r := false
	// 写入数据
	for _, d := range b {
		if p.offset < len(p.Image.Pix) {
			p.Image.Pix[p.offset] = d
		} else {
			p.Image.Pix = append(p.Image.Pix, d)
			r = true
		}
		p.offset++
	}
	// 重置大小
	if r {
		p.Resize()
	}
	return p.offset - n, nil
}

func (p *ImagePack) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	n = p.offset
	for i, _ := range b {
		if p.offset < len(p.Image.Pix) {
			b[i] = p.Image.Pix[p.offset]
			p.offset++
		} else {
			break
		}
	}
	if p.offset-n == 0 {
		return 0, io.EOF
	}
	return p.offset - n, nil
}

func NewPack(size int) *ImagePack {
	w, h := getSize(size)
	return &ImagePack{
		Image:  image.NewNRGBA(image.Rect(0, 0, w, h)),
		offset: 0,
	}
}

func NewDataPack(buf []byte) *ImagePack {
	w, h := getSize(len(buf))
	r := image.Rect(0, 0, w, h)
	return &ImagePack{
		Image:  &image.NRGBA{buf, 4 * w, r},
		offset: 0,
	}
}

func (p *ImagePack) Resize() *image.NRGBA {
	w, h := getSize(len(p.Image.Pix))
	p.Image.Stride = 4 * w
	p.Image.Rect = image.Rect(0, 0, w, h)
	return p.Image
}

func getSize(size int) (w, h int) {
	x := size / 4
	d := math.Sqrt(float64(x))
	return int(math.Ceil(d)), int(math.Ceil(d))
}
