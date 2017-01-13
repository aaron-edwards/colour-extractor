package imgProcess

import (
	"github.com/lucasb-eyer/go-colorful"
)

type Pixel struct {
	R,G,B,A float64
}


func NormRGBtoHSL(col Pixel) (h,s,l float64) {
	h,s,l = colorful.Color{col.R, col.G, col.B}.Hsl()
	return
}


func RgbaToNormPixel(r uint32, g uint32, b uint32, a uint32) (Pixel) {
	return Pixel{ float64(r / 256) / 255.0, float64(g / 256) / 255.0, float64(b / 256) / 255.0, float64(a / 257) / 255.0}
}
