package filters

import (
	"appengine"
	"image"
	"image/color"
	"math"
	"math/rand"
)

func FilterVoronoi(c appengine.Context, m image.Image) image.Image {
	bounds := m.Bounds()
	out := image.NewNRGBA(bounds)
	numClusters := int(math.Sqrt(float64(bounds.Max.Y * bounds.Max.X)))
	// Generates the centroids
	centroids := make(map[int]([]int))
	for i := 0; i < numClusters; i++ {
		centroids[i] = []int{rand.Intn(bounds.Max.X), rand.Intn(bounds.Max.Y)}
	}
	maxval := float64(numClusters * numClusters * numClusters)
	//clSelection := bidimensionalArray(bounds.Max.X, bounds.Max.Y)
	clusterColors := make(map[int]MyColor)

	// Finds the nearest cluster
	clSelection := make([][]int, bounds.Max.Y)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rowSelection := make([]int, bounds.Max.X)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			mindist := maxval
			minCentroid := 0
			for i := 0; i < numClusters; i++ {
				clDistance := distance(centroids[i], x, y)
				if clDistance < mindist {
					mindist = clDistance
					minCentroid = i
				}
			}
			// add the colors to the cluster colors selection
			//r, g, b, a := m.At(x,y).RGBA()
			//c.Infof("Color %v %v %v %v (%v)", r, g, b, a, m.At(x,y))

			rowSelection[x] = minCentroid
			curColor := clusterColors[minCentroid]
			curColor.Add(m.At(x, y))
			clusterColors[minCentroid] = curColor
		}
		clSelection[y] = rowSelection
	}

	// Averages colors
	finalColors := make([]color.Color, numClusters)
	for k, v := range clusterColors {
		//finalColors[k] = m.At(centroids[k][0], centroids[k][1])//colorMean(v)
		finalColors[k] = v.Average()
	}

	// Writes image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			out.Set(x, y, finalColors[clSelection[y][x]])
		}
	}
	return out
}
