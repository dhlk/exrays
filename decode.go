package exrays

import (
	"image"
	"image/color"
	"image/png"
	"io"
)

type Transform func(color.Color) color.Color

var Transforms = map[string]Transform{
	"r": func(rgba color.Color) color.Color {
		b, _, _, _ := rgba.RGBA()
		b = b & 1
		if b == 1 {
			return color.Gray{255}
		}
		return color.Gray{0}
	},
}

func Decode(r io.Reader, w io.Writer, t Transform) error {
	i, err := png.Decode(r)
	if err != nil {
		return err
	}

	b := i.Bounds()
	m := image.NewGray(b)
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			m.Set(x, y, t(i.At(x, y)))
		}
	}

	return png.Encode(w, m)
}
