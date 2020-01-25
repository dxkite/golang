package gogif

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"log"
	"os"
)

func closeIfErr(close func() error) {
	if err := close(); err != nil {
		log.Fatal(err)
	}
}

func MakeGif(input,output string, rect image.Rectangle, delay int)  error {
	f, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer closeIfErr(f.Close)
	w, we := os.Create(output)
	if we != nil {
		return we
	}
	defer closeIfErr(w.Close)

	img, ext, err := image.Decode(f)

	if err != nil {
		return err
	}
	log.Println("read", ext, "image, bounds:", img.Bounds())
	var images []*image.Paletted
	var delays []int
	cp :=  getSubPalette(img)
	for y := 0; y < img.Bounds().Max.Y; y += rect.Dy() {
		pm := image.NewPaletted(rect, cp)
		draw.Draw(pm, rect, img, image.Pt(0, y), draw.Src)
		images = append(images, pm)
		delays = append(delays, delay)
	}
	return gif.EncodeAll(w, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}
