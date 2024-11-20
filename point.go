package wfc

import (
	"errors"
	"fmt"
)

type Point struct {
    X int
    Y int
}

func (p1 *Point) Add(p2 *Point) {
	p1.X += p2.X
	p1.Y += p2.Y
}

func (p1 *Point) Sub(p2 *Point) {
	p1.X -= p2.X
	p1.Y -= p2.Y
}

func (p *Point) Neighbor(direction string) (Point, error) {
	switch direction {
		case "top":
			return Point{p.X, p.Y - 1}, nil
		case "bottom":
			return Point{p.X, p.Y + 1}, nil
		case "left":
			return Point{p.X - 1, p.Y}, nil
		case "right":
			return Point{p.X + 1, p.Y}, nil
		default:
			return *p, errors.New("invalid direction")
	}
}

func (p *Point) OutOfGrid(gridSize int) bool {
	// x is greater than grid size or less than 0
	if p.X >= gridSize || p.X < 0 {
		return true
	}
	// y is greater than grid size or less than 0
	if p.Y >= gridSize || p.Y < 0 {
		return true
	}

	return false
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}