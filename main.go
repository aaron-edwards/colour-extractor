package main

import (
  "os"
  "net/http"
  "github.com/gin-gonic/gin"
  "colour-extractor/routes"
)
import _ "image/png"
import _ "image/jpeg"
import _ "image/gif"

func main() {
  port := os.Getenv("PORT")

  if port == "" {
    port = "8080"
  }

  router := gin.New()
  router.Use(gin.Logger())

  router.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{ "message": "pong" })
  })

  router.GET("/analyse", routes.GetAnalyse)

  router.Run(":" + port)
}
