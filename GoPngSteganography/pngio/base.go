package pngio

import (
	"encoding/binary"
	"image"
	"io"
	"math"
)

type ImagePack struct {
	Image   *image.NRGBA
	Version uint32
	Length  uint32
	offset  int
}

const (
	HeadVersion   = 4
	ContentLength = 4
	HeaderLength  = HeadVersion + ContentLength
)

func (p *ImagePack) Write(b []byte) (n int, err error) {
	n = p.offset
	r := false
	h := HeaderLength
	// 写入数据
	for _, d := range b {
		if p.offset+h < len(p.Image.Pix) {
			p.Image.Pix[h+p.offset] = d
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
	d := p.offset - n
	p.Length += uint32(d)
	p.writeHeader()
	return d, nil
}

func (p *ImagePack) writeHeader() {
	l := make([]byte, 4)
	i := 0
	if len(p.Image.Pix) < HeaderLength {
		p.Image.Pix = make([]byte, HeaderLength)
	}
	binary.BigEndian.PutUint32(l, p.Version)
	p.Image.Pix[i+0] = l[0]
	p.Image.Pix[i+1] = l[1]
	p.Image.Pix[i+2] = l[2]
	p.Image.Pix[i+3] = l[3]
	binary.BigEndian.PutUint32(l, p.Length)
	i += 4
	p.Image.Pix[i+0] = l[0]
	p.Image.Pix[i+1] = l[1]
	p.Image.Pix[i+2] = l[2]
	p.Image.Pix[i+3] = l[3]
}

func (p *ImagePack) readHeader() (version, length uint32) {
	version = binary.BigEndian.Uint32(p.Image.Pix[0:4:4])
	length = binary.BigEndian.Uint32(p.Image.Pix[4:8:8])
	return
}

func (p *ImagePack) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	n = p.offset
	h := HeaderLength
	p.Version, p.Length = p.readHeader()
	for i, _ := range b {
		if uint32(p.offset) < p.Length {
			b[i] = p.Image.Pix[h+p.offset]
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
	img := &ImagePack{
		Image:   image.NewNRGBA(image.Rect(0, 0, w, h)),
		offset:  0,
		Version: 1,
		Length:  0,
	}
	img.writeHeader()
	return img
}

func NewDataPack(buf []byte) *ImagePack {
	w, h := getSize(len(buf))
	r := image.Rect(0, 0, w, h)
	img := &ImagePack{
		Image:   &image.NRGBA{append(make([]byte, HeaderLength), buf...), 4 * w, r},
		Version: 1,
		Length:  uint32(len(buf)),
		offset:  0,
	}
	img.writeHeader()
	return img
}

func (p *ImagePack) Resize() *image.NRGBA {
	w, h := getSize(HeaderLength + int(p.Length))
	p.Image.Stride = 4 * w
	p.Image.Rect = image.Rect(0, 0, w, h)
	return p.Image
}

func getSize(size int) (w, h int) {
	x := size / 4
	d := math.Sqrt(float64(x))
	return int(math.Ceil(d)), int(math.Ceil(d))
}
