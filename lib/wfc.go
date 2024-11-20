package lib

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"maps"
	"math/rand/v2"
	"os"
	"slices"
	"wfc/utils"
)

/* ---------------------------------- CELL ---------------------------------- */

type Cell struct {
	Point *utils.Point
	State []string
	Collapsed bool
}

// collapses cell to a single state
func (c *Cell) observe() error {

	if c.entropy() == -1 {
		return errors.New("found contradiction (no possible states)")
	}

	c.State = []string{c.State[rand.IntN(len(c.State))]}
	c.Collapsed = true;

	return nil
}

// calculates the cells entropy
func (c *Cell) entropy() int {
	return len(c.State) - 1
}

// reduces the possible states to intersection of input
func (c *Cell) narrow(states []interface{}) {
	var possible_states []string

	for _, state := range states {
		for _, possible_state := range c.State {
			if possible_state == state {
				possible_states = append(possible_states, state.(string))
			}
		} 
	} 

	c.State = possible_states
}

// gets the collapsed state of the cell
func (c *Cell) getCollapsedState() string {
	return c.State[0]
}


/* ------------------------- WAVE FUNCTION COLLAPSE ------------------------- */

type Constraints map[string]map[string]interface{}

type WaveFunctionCollapse struct {
	SetName string
	Constraints Constraints
	Grid [][]*Cell
	Resolution int
	Size int
}

// build the grid of the wave function
func (w *WaveFunctionCollapse) Build() {
	states := slices.Collect(maps.Keys(w.Constraints)) // get all constraint keys (possible cell states)

	// allocate and create grid
	w.Grid = make([][]*Cell, w.Size)
	for i := range w.Grid {
		w.Grid[i] = make([]*Cell, w.Size)
		for j := range w.Grid[i] {
			w.Grid[i][j] = &Cell{Point: &utils.Point{X: i, Y: j}, State: states, Collapsed: false}
		}
	}
}

// reset the wave function (calls build)
func (w *WaveFunctionCollapse) Reset() {
	w.Build()
}

// find the minimum entropy cell of the wave function
func (w *WaveFunctionCollapse) findMinEntropyCell() *Cell {
	var minEntropyCell *Cell = nil

	for i := range w.Grid {
		for _, cell := range w.Grid[i] {
			if !cell.Collapsed && (minEntropyCell == nil || cell.entropy() < minEntropyCell.entropy()) {
				minEntropyCell = cell
			}
		}
	}

	return minEntropyCell
}

// collapse the minimum entropy cell and propagate states
func (w *WaveFunctionCollapse) CollapseNext() bool {
	cell := w.findMinEntropyCell()

	// no min entropy cell (all cells collapsed) -> complete
	if cell == nil {
		return true
	}

	cell.observe() // observe and collapse current cell state

	for _, direction := range [...]string{"top", "bottom", "left", "right"} {
		neighboringPoint, _ := cell.Point.Neighbor(direction)

		// neighbor is out of grid -> move on
		if neighboringPoint.OutOfGrid(w.Size) {
			continue
		}

		neighbor := w.Grid[neighboringPoint.X][neighboringPoint.Y]

		// neighbor is already collapsed -> move on
		if neighbor.Collapsed {
			continue
		}

		// neighbor within grid and not collapsed -> update possible states
		possibilities := w.Constraints[cell.getCollapsedState()][direction]
		neighbor.narrow(possibilities.([]interface{}))
	}

	return false
}


// collapse all cells
func (w *WaveFunctionCollapse) Collapse() error {
	for {
		// all cells collapsed -> exit
		if w.CollapseNext() {
			break
		}
	}

	return nil
}

// save wave to image
func (w *WaveFunctionCollapse) Save() error {
	// new image
	s := w.Size * w.Resolution
	img := image.NewRGBA(image.Rect(0, 0, s, s))

	// overlay cells
	for _, row := range w.Grid {
		for _, cell := range row {

			// get the collapsed state image
			file, err := os.Open(fmt.Sprintf("assets/%v/%v.png", w.SetName, cell.getCollapsedState()))
			if err != nil {
				log.Fatal(err)
			}

			asset, _, err := image.Decode(file)
			if err != nil {
				log.Fatal(err)
			}

			// draw state
			draw.Draw(img, asset.Bounds().Add(image.Point{X: cell.Point.X * w.Resolution, Y: cell.Point.Y * w.Resolution}), asset, image.Point{0, 0}, draw.Over)

			file.Close()
		}
	}

	// save image
	outputFile, err := os.Create("compiled.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	png.Encode(outputFile, img)

	return nil
}