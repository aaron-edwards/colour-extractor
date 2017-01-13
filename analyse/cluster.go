package analyse

import(
  "image"
  "time"
  "log"
	"github.com/mdesenfants/gokmeans"
	"colour-extractor/img-process"
	"github.com/lucasb-eyer/go-colorful"
)

func toNode(data [][3]float64) ([]gokmeans.Node) {
	var nodes []gokmeans.Node
	for x := 0; x < len(data); x++ {
		nodes = append(nodes, gokmeans.Node{data[x][0], data[x][1], data[x][2]})
	}

	return nodes
}

func toHex(col gokmeans.Node) string {
  return colorful.Hsl(col[0], col[1], col[2]).Hex()
}

//func toRGB(col gokmeans.Node) (r, g, b, a uint8) {
//  return colorful.Hsl(col[0], col[1], col[2]).RGB255()
//}

type HSL struct {
  H float64
  S float64
  L float64
}

type ClusterData struct {
  CenterHex string
  CenterHSL HSL
  PixelRatio float64
}

func groupNodes(centroids []gokmeans.Node, data []gokmeans.Node) ([]ClusterData) {
  var colourGroups = make(map[int][]gokmeans.Node)
  for _, node := range data {
    index := gokmeans.Nearest(node, centroids)
    colourGroups[index] = append(colourGroups[index], node)
  }

  var clusters  []ClusterData
  for index, group := range colourGroups {
    centroid := centroids[index]
    hsl := HSL{ H: centroid[0], S: centroid[1], L: centroid[2]}
    hex := toHex(centroid)
   // rgb := toRGB(centroid)
    ratio := float64(len(group)) / float64(len(data))
    clusters = append(clusters, ClusterData{ CenterHex: hex, CenterHSL: hsl, PixelRatio: ratio })
  }

  return clusters
}

func Cluster(img image.Image) ([]ClusterData, []gokmeans.Node) {
	pixels := imgProcess.GetPixels(img, 0.75)
	nodes := toNode(pixels)

  start := time.Now()
  k := 5
	_, centroids := gokmeans.Train(nodes, k, 20)
  log.Printf("time to cluster:\t%s", time.Since(start))

  groups := groupNodes(centroids, nodes)

  return groups, centroids
}


