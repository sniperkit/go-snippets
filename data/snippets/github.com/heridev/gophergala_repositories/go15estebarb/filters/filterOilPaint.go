package filters

import (
	"appengine"
	"image"
)

func FilterOilPaint(c appengine.Context, m image.Image) image.Image {
	bounds := m.Bounds()
	out := image.NewNRGBA(bounds)
	ys := bounds.Max.Y
	xs := bounds.Max.X
	radius := 5
	intensityLevels := 20

	intensityMap := make([][]uint8, ys)
	for y := 0; y < ys; y++ {
		intensityRow := make([]uint8, xs)
		for x := 0; x < xs; x++ {
			currentColor := m.At(x, y)
			r, g, b, _ := currentColor.RGBA()
			//c.Infof("Color %v %v %v (%v)", r, g, b, currentColor)
			ci := uint8(int(r+g+b) / 3.0 * intensityLevels / 255.0 / 255.0)
			intensityRow[x] = ci
		}
		intensityMap[y] = intensityRow
	}

	for y := 0; y < ys; y++ {
		for x := 0; x < xs; x++ {
			intensities := make([]MyColor, intensityLevels+1)
			for y2 := IntMax(0, y-radius); y2 < IntMin(ys, y+radius); y2++ {
				for x2 := IntMax(0, x-radius); x2 < IntMin(xs, x+radius); x2++ {
					currentColor := m.At(x2, y2)
					//r,g,b,_ := currentColor.RGBA()
					//c.Infof("Color %v %v %v (%v)", r, g, b, currentColor)
					//ci := int(int(r+g+b)/3.0*intensityLevels/255.0/255.0)
					ci := intensityMap[y2][x2]
					//c.Infof("Intensities %v of %v", ci, len(intensities))
					newColor := intensities[ci]
					newColor.Add(currentColor)
					intensities[ci] = newColor
				}
			}
			newColor := intensities[0]
			for _, v := range intensities {
				if newColor.C < v.C {
					newColor = v
				}
			}
			out.Set(x, y, newColor.Average())
		}
	}
	return out
}
