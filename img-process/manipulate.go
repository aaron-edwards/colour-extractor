package imgProcess

import (
  "math"
  "image"
  "github.com/nfnt/resize"
)

func round(f float64) (uint) {
  return uint(math.Floor(f + 0.5))
}

func imageSize(img image.Image) (int, int) {
  bounds := img.Bounds()
  width := bounds.Max.X - bounds.Min.X
  height := bounds.Max.Y - bounds.Min.Y

  return width, height
}

func ResizeImage(img image.Image, pixels int) (image.Image) {
  width, height := imageSize(img)
  ratio := float64(width) / float64(height)
  newWidth := round(math.Sqrt(float64(pixels) * ratio))

  return resize.Resize(newWidth, 0, img, resize.Bilinear)
}

func GetPixels(img image.Image) ([] RGBA) {
  width, height := imageSize(img)

  var pixels = make([]RGBA, width * height)
  for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
      pixels[x + (y * width)] = rgbaToPixel(img.At(x, y).RGBA())
    }
  }

  return pixels
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) RGBA {
    return RGBA{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

type RGBA struct {
	R,G,B,A int
}
