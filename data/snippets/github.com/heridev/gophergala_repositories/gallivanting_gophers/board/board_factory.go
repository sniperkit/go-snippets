package board

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gophergala/gallivanting_gophers/pieces"
)

// Factory encapsulates generation of game boards and testing.
type Factory struct {
	random *rand.Rand
}

// NewFactory creates a new factory.
func NewFactory() *Factory {
	fmt.Printf("Creating Board Factory...")
	f := &Factory{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	return f
}

// GetBoard currently returns a standard board. Should be extended to support
// multiple board types.
func (g *Factory) GetBoard() *FixBoard {
	b := g.newStandardBoard()

	return b
}

// Creates a standard board setup.
func (g *Factory) newStandardBoard() *FixBoard {
	var b *FixBoard
	var valid = false

	for !valid {
		b = newEmptyBoard()

		// Middle Blocks
		b.SetWalls(119, true, true, false, false)
		b.SetWalls(120, false, true, true, false)
		b.SetWalls(135, true, false, false, true)
		b.SetWalls(136, false, false, true, true)

		g.generateOuterWalls(b)
		g.generateShapes(b)

		// TODO: Board Validation (remove invalid boards)

		if !g.generateGoals(b) {
			fmt.Println("Invalid board generated due to Target locations! Regenerating...")
			continue
		}

		if !g.generateGophers(b) {
			fmt.Println("Invalid board generated due to Gopher locations! Regenerating...")
			continue
		}

		valid = true
	}

	return b
}

// Generates the outer walls.
func (g *Factory) generateOuterWalls(b *FixBoard) {
	topFirst := 1 + g.random.Intn(5)
	topSecond := 9 + g.random.Intn(5)
	b.SetWall(topFirst, DirRight)
	b.SetWall(topSecond, DirRight)

	leftTop := 1 + g.random.Intn(5)
	leftBottom := 9 + g.random.Intn(5)

	b.SetWall(leftTop*16, DirBottom)
	b.SetWall(leftBottom*16, DirBottom)

	rightTop := 1 + g.random.Intn(5)
	rightBottom := 9 + g.random.Intn(5)

	b.SetWall((rightTop+1)*16-1, DirBottom)
	b.SetWall((rightBottom+1)*16-1, DirBottom)

	bottomFirst := 1 + g.random.Intn(5)
	bottomSecond := 9 + g.random.Intn(5)

	b.SetWall(16*15+bottomFirst, DirRight)
	b.SetWall(16*15+bottomSecond, DirRight)
}

// Generates the angle shapes on the map.
func (g *Factory) generateShapes(b *FixBoard) {
	duplicates := make(map[int]struct{})

	// Temporary fix for middle. This needs to be changed to support non-standard boards.
	duplicates[119] = struct{}{}
	duplicates[120] = struct{}{}
	duplicates[135] = struct{}{}
	duplicates[136] = struct{}{}

	for x := 1; x < 15; x++ {
		multiple := rand.Intn(10) == 9
		count := 1
		if multiple {
			count = 2
		}
		for ; count > 0; count-- {
			square := 1 + rand.Intn(14)

			var squareID = square + x*16

			if _, exists := duplicates[squareID]; exists {
				continue
			}

			duplicates[squareID] = struct{}{}

			up := rand.Intn(2)
			right := rand.Intn(2)

			b.SetWalls(squareID, right == 0, up == 1, right == 1, up == 0)

			b.locations = append(b.locations, squareID)
		}
	}
}

// Sets up the goal locations on the board.
func (g *Factory) generateGoals(b *FixBoard) bool {
	duplicates := make(map[int]struct{})

	for x := 0; x < 5; x++ {
		var idx int
		unique := false

		// Prevents infinite loops for maps which cannot support enough locations
		if len(duplicates) == len(b.locations) {
			return false
		}

		for !unique {
			idx = rand.Intn(len(b.locations))
			if _, exists := duplicates[b.locations[idx]]; !exists {
				duplicates[b.locations[idx]] = struct{}{}
				unique = true
			}
		}
		b.Goals = append(b.Goals, pieces.Goal{
			Type:     pieces.GoalFromInt(x),
			Location: b.locations[idx],
		})
	}

	return true
}

// Sets up the gopher starting locations on the board.
func (g *Factory) generateGophers(b *FixBoard) bool {
	duplicates := make(map[int]struct{})

	// Temporary fix for middle. This needs to be changed to support non-standard boards.
	duplicates[119] = struct{}{}
	duplicates[120] = struct{}{}
	duplicates[135] = struct{}{}
	duplicates[136] = struct{}{}

	// Pre-load target locations, do not generate new boards with gophers already on targets.
	for x := 0; x < len(b.Goals); x++ {
		duplicates[b.Goals[x].Location] = struct{}{}
	}

	for x := 0; x < 5; x++ {
		var idx int
		unique := false

		for !unique {
			idx = rand.Intn(256)
			if _, exists := duplicates[idx]; !exists {
				duplicates[idx] = struct{}{}
				unique = true
			}
		}

		b.Gophers = append(b.Gophers, &pieces.Gopher{
			Type:     pieces.GopherFromInt(x),
			Location: idx,
		})
	}

	return true
}
