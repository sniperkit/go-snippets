package algorithm

import (
	"github.com/gophergala/edrans-smartcity/models"

	"fmt"
)

func getCandidates(c *models.City, origin, dest *models.Node, visited []int) []models.Path {
	if c.LastError != nil {
		return nil
	}
	vlen := len(visited)
	var paths = make([]models.Path, 0)
	for i := 0; i < len(origin.Outputs); i++ {
		if origin.Outputs[i].DestinyID == dest.ID {
			paths = append(paths, models.Path{Reached: true, Links: []models.Link{origin.Outputs[i]}})
			continue
		}
		if alreadyVisited(origin.Outputs[i].DestinyID, visited) {
			continue
		}
		visited = append(visited, origin.Outputs[i].DestinyID)
		subPaths := getCandidates(c, c.GetNode(origin.Outputs[i].DestinyID), dest, visited)
		for j := 0; j < len(subPaths); j++ {
			lnks := subPaths[j].Links
			lnks = append(lnks, origin.Outputs[i])
			paths = append(paths, models.Path{Links: lnks, Reached: subPaths[j].Reached})
		}
	}
	if vlen == 0 && len(paths) == 0 {
		c.LastError = fmt.Errorf("There's no way to the requested address")
	}
	return paths
}

func alreadyVisited(ID int, visited []int) bool {
	for i := 0; i < len(visited); i++ {
		if visited[i] == ID {
			return true
		}
	}
	return false
}

func GetPaths(c *models.City, origin, destiny int) ([]models.Path, error) {
	if origin == destiny {
		return nil, fmt.Errorf("Already at destiny")
	}
	org := c.GetNode(origin)
	dest := c.GetNode(destiny)
	candidates := getCandidates(c, org, dest, nil)
	return sortLinks(candidates), c.LastError
}

func sortLinks(paths []models.Path) []models.Path {
	for i := 0; i < len(paths); i++ {
		var lnk = make([]models.Link, len(paths[i].Links))
		for j := 0; j < len(lnk); j++ {
			lnk[j] = paths[i].Links[len(lnk)-1-j]
		}
		paths[i].Links = lnk
	}
	return paths
}

func ChooseBest(paths []models.Path) models.Path {
	var better = paths[0]
	for i := 1; i < len(paths); i++ {
		if paths[i].Estimate < better.Estimate {
			better = paths[i]
		}
	}
	return better
}

func CalcEstimatesForVehicle(v *models.Vehicle, paths []models.Path) []models.Path {
	for i := 0; i < len(paths); i++ {
		if len(paths[i].Links) == 0 {
			continue
		}
		paths[i].Weights = make([]int, 0)
		weight := paths[i].Links[0].Weight
		paths[i].Weights = append(paths[i].Weights, paths[i].Links[0].Weight)
		for j := 1; j < len(paths[i].Links); j++ {

			weight += paths[i].Links[j].Weight
			paths[i].Weights = append(paths[i].Weights, paths[i].Links[j].Weight)
		}
		paths[i].Estimate = weight
	}
	return paths
}
