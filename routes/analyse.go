package routes

import (
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

	start := time.Now()

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

	log.Printf("Downloading took %s", time.Since(start))
  downloadTime := time.Now()

	img = imgProcess.ResizeImage(img, 50 * 50)
	pixels := imgProcess.GetPixels(img, 0.75)
	nodes := toNode(pixels)

	log.Printf("Processing took %s", time.Since(downloadTime))
  imgProcessTime := time.Now()

  k := 5
	_, centroids := gokmeans.Train(nodes, k, 20)

	log.Printf("Clustering took %s", time.Since(imgProcessTime))
  clusterTime := time.Now()


  var colourGroups = make(map[int]int)
  for _, node := range nodes {
    index := gokmeans.Nearest(node, centroids)
    colourGroups[index] ++
  }

  var colours = make(map[string]float64)
  for index, centroid := range centroids {
    colours[toHex(centroid)] = float64(colourGroups[index]) / float64(len(nodes))
  }


	log.Printf("Grouping took %s", time.Since(clusterTime))

	c.JSON(http.StatusOK, gin.H{ "imageUrl": imageUrl, "colours": colours})
}

func toHex(col gokmeans.Node) string {
  return colorful.Hsl(col[0], col[1], col[2]).Hex()
}

