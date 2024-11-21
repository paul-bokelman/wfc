package wfc

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"maps"
	"math/rand/v2"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

/* ---------------------------------- CELL ---------------------------------- */

type Cell struct {
	Point *Point
	State []string
	Collapsed bool
}

// collapses cell to a single state
func (this *Cell) observe() error {

	if this.entropy() == -1 {
		return errors.New("found contradiction (no possible states)")
	}

	this.State = []string{this.State[rand.IntN(len(this.State))]}
	this.Collapsed = true;

	return nil
}

// calculates the cells entropy
func (this *Cell) entropy() int {
	return len(this.State) - 1
}

// reduces the possible states to intersection of input
func (this *Cell) narrow(states []interface{}) {
	var possible_states []string

	for _, state := range states {
		for _, possible_state := range this.State {
			if possible_state == state {
				possible_states = append(possible_states, state.(string))
			}
		} 
	} 

	this.State = possible_states
}

// gets the collapsed state of the cell
func (this *Cell) getCollapsedState() string {
	return this.State[0]
}


/* ------------------------- WAVE FUNCTION COLLAPSE ------------------------- */

type Rules map[string]map[string]interface{}

type ConstraintsSet struct {
	Resolution int      `yaml:"resolution"`
	Rules      Rules  	`yaml:"rules"`
}

type WaveFunctionCollapseConfig struct {
	AssetsPath string
	OutputPath string 
	GridSize int
}

type WaveFunctionCollapse struct {
	assetsPath string
	outputPath string
	size int 
	resolution int
	rules Rules
	grid [][]*Cell
}

// initialize the wave function given an asset path and grid size
func (this *WaveFunctionCollapse) Initialize(config WaveFunctionCollapseConfig) error {
    data, err := os.ReadFile(fmt.Sprintf("%v/constraints.yaml", config.AssetsPath))

	if err != nil {
		return errors.New("failed to read file")
	}

	constraints := ConstraintsSet{}
    
	err = yaml.Unmarshal([]byte(data), &constraints)

	if err != nil {
		return errors.New("failed to parse yaml file")
	}

	// assign local vars 
	this.assetsPath = config.AssetsPath
	this.outputPath = config.OutputPath
	this.size = config.GridSize
	this.resolution = constraints.Resolution
	this.rules = constraints.Rules
	
	this.build() // build the wfc grid 
	return nil
}

// build the grid of the wave function
func (this *WaveFunctionCollapse) build() {
	states := slices.Collect(maps.Keys(this.rules)) // get all constraint keys (possible cell states)

	// allocate and create grid
	this.grid = make([][]*Cell, this.size)
	for i := range this.grid {
		this.grid[i] = make([]*Cell, this.size)
		for j := range this.grid[i] {
			this.grid[i][j] = &Cell{Point: &Point{X: i, Y: j}, State: states, Collapsed: false}
		}
	}
}

// reset the wave function (calls build)
func (this *WaveFunctionCollapse) Reset() {
	this.build()
}

// find the minimum entropy cell of the wave function
func (this *WaveFunctionCollapse) findMinEntropyCell() *Cell {
	var minEntropyCell *Cell = nil

	for i := range this.grid {
		for _, cell := range this.grid[i] {
			if !cell.Collapsed && (minEntropyCell == nil || cell.entropy() < minEntropyCell.entropy()) {
				minEntropyCell = cell
			}
		}
	}

	return minEntropyCell
}

// collapse the minimum entropy cell and propagate states
func (this *WaveFunctionCollapse) CollapseNext() bool {
	cell := this.findMinEntropyCell()

	// no min entropy cell (all cells collapsed) -> complete
	if cell == nil {
		return true
	}

	cell.observe() // observe and collapse current cell state

	for _, direction := range [...]string{"top", "bottom", "left", "right"} {
		neighboringPoint, _ := cell.Point.Neighbor(direction)

		// neighbor is out of grid -> move on
		if neighboringPoint.OutOfGrid(this.size) {
			continue
		}

		neighbor := this.grid[neighboringPoint.X][neighboringPoint.Y]

		// neighbor is already collapsed -> move on
		if neighbor.Collapsed {
			continue
		}

		// neighbor within grid and not collapsed -> update possible states
		possibilities := this.rules[cell.getCollapsedState()][direction]
		neighbor.narrow(possibilities.([]interface{}))
	}

	return false
}


// collapse all cells
func (this *WaveFunctionCollapse) Collapse() error {
	for {
		// all cells collapsed -> exit
		if this.CollapseNext() {
			break
		}
	}

	return nil
}

// save wave to image
func (this *WaveFunctionCollapse) Save() error {
	// new image
	s := this.size * this.resolution
	img := image.NewRGBA(image.Rect(0, 0, s, s))

	// overlay cells
	for _, row := range this.grid {
		for _, cell := range row {

			// get the collapsed state image
			file, err := os.Open(fmt.Sprintf("%v/%v.png", this.assetsPath, cell.getCollapsedState()))

			if err != nil {
				return errors.New("failed to get asset file")
			}

			asset, _, err := image.Decode(file)

			if err != nil {
				return errors.New("failed to decode asset file")
			}

			// draw state
			draw.Draw(img, asset.Bounds().Add(image.Point{X: cell.Point.X * this.resolution, Y: cell.Point.Y * this.resolution}), asset, image.Point{0, 0}, draw.Over)

			file.Close()
		}
	}

	// save image
	outputFile, err := os.Create(this.outputPath)

	if err != nil {
		return errors.New("failed to create file")
	}

	defer outputFile.Close()
	
	png.Encode(outputFile, img)
	
	return nil
}