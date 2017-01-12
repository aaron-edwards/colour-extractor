package routes

import (
  "log"
  "time"
  "net/http"
  "github.com/gin-gonic/gin"
  "colour-extractor/img-process"
  "image"
  _ "image/png"
  _ "image/jpeg"
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

  start := time.Now()

  img = imgProcess.ResizeImage(img, 100 * 100)
  imgProcess.GetPixels(img)

  elapsed := time.Since(start)
  log.Printf("Binomial took %s", elapsed)

  c.JSON(http.StatusOK, gin.H{ "imageUrl": imageUrl })
}
