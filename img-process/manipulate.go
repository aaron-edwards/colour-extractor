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

func GetPixels(img image.Image, alphaLimit float64) ([][3]float64) {
  width, height := imageSize(img)

  var pixels [][3]float64
  for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
			pixel := RgbaToNormPixel(img.At(x,y).RGBA())
			if float64(pixel.A) >= alphaLimit {
				h,c,l := NormRGBtoHSL(pixel)
      	pixels = append(pixels, [3]float64{h,c,l})
			}
    }
  }
  return pixels
}


