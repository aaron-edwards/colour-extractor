package routes

import (
  "log"
  "bytes"
  "strconv"
  "net/http"
  "github.com/gin-gonic/gin"
  "colour-extractor/img-process"
  "colour-extractor/analyse"
  "image"
  "image/png"
  _ "image/jpeg"
  "image/color"
  "image/draw"
)


func GetAnalyse(c *gin.Context) {
  imageUrl := c.Query("imageUrl")

  if imageUrl == "" {
    c.JSON(http.StatusBadRequest, gin.H{ "error": "missing 'imageUrl' query param" })
    return
  }

  response, err := http.Get(imageUrl)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{ "error": err })
    return
  }
  defer response.Body.Close()

  img, _, err := image.Decode(response.Body)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{ "error": err })
  }

  resizedImage := imgProcess.ResizeImage(img, 100 * 100)
  raw, clusters := analyse.Cluster(resizedImage)

  if c.Query("image") == "true" {
    //imgProcess.generateImage(img, raw, clusters)
   returnImage := imgProcess.ResizeImage(img, 500 * 500)

   if bigImage, ok := returnImage.(*image.RGBA); ok {
    rect := image.Rect(0,0,50,50)
    col := color.RGBA{0, 0, 255, 255}
    draw.Draw(bigImage, rect, &image.Uniform{col}, image.ZP, draw.Src)

    buffer := new(bytes.Buffer)
    if err := png.Encode(buffer, returnImage); err != nil {
      log.Println("unable to encode image.")
    }

    w := c.Writer

    w.Header().Set("Content-Type", "image/png")
    w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
    if _, err := w.Write(buffer.Bytes()); err != nil {
      log.Println("unable to write image.")
    }
  }

  } else {
    c.JSON(http.StatusOK, gin.H{ "clusters": clusters, "raw" : raw })
  }
}


