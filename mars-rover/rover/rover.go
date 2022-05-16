package rover

import "errors"

var (
	ErrInvalidDirection = errors.New("invalid direction")
	ErrNotDropped       = errors.New("rover not dropped")
)

type Direction string

const (
	North Direction = "N"
	East  Direction = "E"
	South Direction = "S"
	West  Direction = "W"
)

type Rover struct {
	pos *Position
}

type Position struct {
	X   int
	Y   int
	Dir Direction
}

func (r *Rover) Drop(x, y int, direction Direction) (*Position, error) {
	if !isValidDirection(direction) {
		return nil, ErrInvalidDirection
	}

	r.pos = &Position{x, y, direction}
	return r.pos, nil
}

func (r *Rover) Move(instructions string) (*Position, error) {
	if r.pos == nil {
		return nil, ErrNotDropped
	}

	return r.pos, nil
}

func isValidDirection(direction Direction) bool {
	return direction == North || direction == East || direction == West || direction == South
}
