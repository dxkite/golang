package gogif

import (
	"image"
	"image/color"
	"image/color/palette"
)

func inPalette(p color.Palette, c color.Color) int {
	ret := -1
	for i, v := range p {
		if v == c {
			return i
		}
	}
	return ret
}

func getSubPalette(m image.Image) color.Palette {
	p := color.Palette{color.RGBA{0x00,0x00,0x00,0x00}}
	p9 := color.Palette(palette.Plan9)
	b := m.Bounds()
	black := false
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := m.At(x, y)
			cc := p9.Convert(c)
			if cc == p9[0] {
				black = true
			}
			if inPalette(p, cc) == -1 {
				p = append(p, cc)
			}
		}
	}
	if len(p) < 256 && black == true {
		p[0] = color.RGBA{0x00,0x00,0x00,0x00} // transparent
		p = append(p, p9[0])
	}
	return p
}

