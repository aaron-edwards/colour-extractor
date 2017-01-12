package routes

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"colour-extractor/img-process"
	"github.com/mdesenfants/gokmeans"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	_ "image/png"
	_ "image/jpeg"
)

func toNode(data [][3]float64) ([]gokmeans.Node) {
	var nodes []gokmeans.Node
	for x := 0; x < len(data); x++ {
		nodes = append(nodes, gokmeans.Node{data[x][0], data[x][1], data[x][2]})
	}

	return nodes
}

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
	pixels := imgProcess.GetPixels(img, 0.75)
	nodes := toNode(pixels)

	var colours []string
	_, centroids := gokmeans.Train(nodes, 5, 50)

		for _, centroid := range centroids {
			colours = append(colours, toHex(centroid))
		}

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)


		for _, nodes := range nodes {
			index := gokmeans.Nearest(nodes, centroids)
			fmt.Println(toHex(nodes), nodes, "belongs in cluster", index+1, ".")
		}

	c.JSON(http.StatusOK, gin.H{ "imageUrl": imageUrl, "colours": colours, "elapsed" : elapsed})
}

func toHex(rgb gokmeans.Node) string {
	return colorful.Color{rgb[0], rgb[1], rgb[2]}.Hex()
}

func labToHex(lab gokmeans.Node) string {
	return colorful.Lab(lab[0], lab[1], lab[2]).Hex()
}
