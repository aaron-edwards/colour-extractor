package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
    port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{ "message": "pong" })
	})

	router.Run(":" + port)
}
