package factory

import (
	"github.com/gophergala/edrans-smartcity/models"

	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	MIN_WEIGHT_AMBULANCE         = 10
	MIN_WEIGHT_FIREFIGHT_VEHICLE = 15
	MIN_WEIGHT_POLICE_CARS       = 5
	MAX_POLICE_CARS              = 5
	MAX_FIREFIGHT_VEHICLES       = 5
	MAX_AMBULANCES               = 5
)

func init() {
	rand.Seed(time.Now().Unix())
}

func CreateRectangularCity(height int, width int, name string) (myCity *models.City, err error) {
	if height <= 2 || width <= 2 {
		return nil, fmt.Errorf("City size must be greater than 2")
	}

	numNodes := height * width
	var city = make([]models.Node, numNodes)

	// ID are 1-based and arrays are 0-based
	currID := 1
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			neighbourList := getNeighbours(i, j, height, width)

			linkList := make([]models.Link, 0)
			for k := 0; k < len(neighbourList); k++ {
				if currNeighbour := neighbourList[k]; currNeighbour != -1 {
					linkList = append(linkList,
						models.Link{
							Name:      strconv.Itoa(currID) + "-to-" + strconv.Itoa(currNeighbour) + "St.",
							OriginID:  currID,
							DestinyID: currNeighbour,
							Weight:    getRandomWeight(),
						},
					)
				}
			}

			city[currID-1] = models.Node{
				ID:       currID,
				Location: []int{j, -i},
				Outputs:  linkList,
			}
			currID++
		}
	}

	myCity, err = models.NewCity(city, name, height, width)
	if err != nil {
		return nil, err
	}

	myCity.AddService(models.SERVICE_HOSPITAL, models.NewPublicServicePosition(myCity, numNodes), rand.Intn(MAX_POLICE_CARS)+1, MIN_WEIGHT_AMBULANCE)
	myCity.AddService(models.SERVICE_FIREFIGHTER, models.NewPublicServicePosition(myCity, numNodes), rand.Intn(MAX_FIREFIGHT_VEHICLES)+1, MIN_WEIGHT_FIREFIGHT_VEHICLE)
	myCity.AddService(models.SERVICE_POLICE, models.NewPublicServicePosition(myCity, numNodes), rand.Intn(MAX_POLICE_CARS)+1, MIN_WEIGHT_POLICE_CARS)
	myCity.LaunchVehicles()
	return
}

func getNeighbours(i, j, m, n int) (list []int) {
	list = make([]int, 0)
	vert := calcVerticalNeighbourID(i, j, m, n)
	if vert != -1 {
		list = append(list, vert)
	}
	horiz := calcHorizontalNeighbourID(i, j, m, n)
	if horiz != -1 {
		list = append(list, horiz)
	}

	return
}

func calcHorizontalNeighbourID(i, j, m, n int) int {
	var output int

	if i%2 == 0 {
		// then strets go right
		if j+1 == n {
			return -1
		}
		output = i*n + j + 2
	} else {
		// then streets go left
		if j == 0 {
			return -1
		}
		output = i*n + j
	}

	return output
}

func calcVerticalNeighbourID(i, j, m, n int) int {
	var output int

	if j%2 == 0 {
		// then streets go down
		if i+1 == m {
			return -1
		}
		output = n*(i+1) + j + 1
	} else {
		// then streets go up
		if i == 0 {
			return -1
		}
		output = n*(i-1) + j + 1
	}

	return output
}

func getRandomNode(numNodes int) int {
	return rand.Intn(numNodes) + 1 // return [1, n] random
}

func getRandomWeight() int {
	return rand.Intn(100)
}

func SampleCity() *models.City {
	var city = make([]models.Node, 16)

	city[0] = models.Node{ID: 1, Outputs: []models.Link{models.Link{Name: "Roca", OriginID: 1, DestinyID: 2, Weight: 30}, models.Link{Name: "Pellegrini", OriginID: 1, DestinyID: 5, Weight: 30}}}
	city[1] = models.Node{ID: 2, Outputs: []models.Link{models.Link{Name: "Roca", OriginID: 2, DestinyID: 3, Weight: 30}}}
	city[2] = models.Node{ID: 3, Outputs: []models.Link{models.Link{Name: "Roca", OriginID: 3, DestinyID: 4, Weight: 30}, models.Link{Name: "Irigoyen", OriginID: 3, DestinyID: 7, Weight: 35}}}
	city[3] = models.Node{ID: 4, Outputs: []models.Link{}}
	city[4] = models.Node{ID: 5, Outputs: []models.Link{models.Link{Name: "Pellegrini", OriginID: 5, DestinyID: 9, Weight: 30}}}
	city[5] = models.Node{ID: 6, Outputs: []models.Link{models.Link{Name: "Rivadavia", OriginID: 6, DestinyID: 5, Weight: 35}, models.Link{Name: "Irigoyen", OriginID: 6, DestinyID: 2, Weight: 35}}}
	city[6] = models.Node{ID: 7, Outputs: []models.Link{models.Link{Name: "Rivadavia", OriginID: 7, DestinyID: 6, Weight: 45}, models.Link{Name: "Palacios", OriginID: 7, DestinyID: 11, Weight: 45}}}
	city[7] = models.Node{ID: 8, Outputs: []models.Link{models.Link{Name: "Rivadavia", OriginID: 8, DestinyID: 7, Weight: 35}, models.Link{Name: "Justo", OriginID: 8, DestinyID: 12, Weight: 30}}}
	city[8] = models.Node{ID: 9, Outputs: []models.Link{models.Link{Name: "Mitre", OriginID: 9, DestinyID: 10, Weight: 35}, models.Link{Name: "Pellegrini", OriginID: 9, DestinyID: 13, Weight: 30}}}
	city[9] = models.Node{ID: 10, Outputs: []models.Link{models.Link{Name: "Irigoyen", OriginID: 10, DestinyID: 6, Weight: 45}, models.Link{Name: "Mitre", OriginID: 10, DestinyID: 11, Weight: 45}}}
	city[10] = models.Node{ID: 11, Outputs: []models.Link{models.Link{Name: "Palacios", OriginID: 11, DestinyID: 15, Weight: 35}, models.Link{Name: "Mitre", OriginID: 11, DestinyID: 12, Weight: 35}}}
	city[11] = models.Node{ID: 12, Outputs: []models.Link{models.Link{Name: "Justo", OriginID: 12, DestinyID: 8, Weight: 30}}}
	city[12] = models.Node{ID: 13, Outputs: []models.Link{}}
	city[13] = models.Node{ID: 14, Outputs: []models.Link{models.Link{Name: "Irigoyen", OriginID: 14, DestinyID: 10, Weight: 35}, models.Link{Name: "Urquiza", OriginID: 14, DestinyID: 13, Weight: 30}}}
	city[14] = models.Node{ID: 15, Outputs: []models.Link{models.Link{Name: "Urquiza", OriginID: 15, DestinyID: 14, Weight: 30}}}
	city[15] = models.Node{ID: 16, Outputs: []models.Link{models.Link{Name: "Justo", OriginID: 16, DestinyID: 12, Weight: 30}, models.Link{Name: "Urquiza", OriginID: 16, DestinyID: 15, Weight: 30}}}
	myCity, _ := models.NewCity(city, "Fake Buenos Aires", 4, 4)

	myCity.AddService(models.SERVICE_HOSPITAL, 10, 5, 10)
	myCity.AddService(models.SERVICE_FIREFIGHTER, 11, 5, 15)
	myCity.AddService(models.SERVICE_POLICE, 16, 5, 5)

	myCity.LaunchVehicles()
	return myCity
}

/*

d      u      d      u

1 ---- 2 ---- 3 ---- 4    r Roca
|      |      |      |
5 ---- 6 ---- 7 ---- 8    l Rivadavia
|      |      |      |
9 ---- a ---- b ---- c    r Mitre
|      |      |      |
d ---- e ---- f ---- g    l Urquiza

- Pellegrini
- Irigoyen
- Palacios
- Justo

*/
