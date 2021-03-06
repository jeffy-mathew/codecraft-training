package rover

import (
	"errors"
)

var (
	ErrNotDropped = errors.New("rover not dropped")
)

const (
	NorthNumber = iota
	EastNumber
	SouthNumber
	WestNumber
)

type Direction interface {
	GetNumber() int
}

type North struct{}

func (n North) GetNumber() int {
	return NorthNumber
}

type East struct{}

func (e East) GetNumber() int {
	return EastNumber
}

type South struct{}

func (s South) GetNumber() int {
	return SouthNumber
}

type West struct{}

func (w West) GetNumber() int {
	return WestNumber
}

func getDirectionByNumber(currDir int) Direction {
	switch currDir {
	case 0:
		return North{}
	case 1:
		return East{}
	case 2:
		return South{}
	case 3:
		return West{}
	}
	return nil
}

type Rover struct {
	pos *Position
}

type Position struct {
	X   int
	Y   int
	Dir Direction
}

func (r *Rover) Drop(x, y int, direction Direction) (*Position, error) {

	r.pos = &Position{x, y, direction}
	return r.pos, nil
}

func (r *Rover) Move(instructions string) (*Position, error) {
	if r.pos == nil {
		return nil, ErrNotDropped
	}

	for _, instruction := range instructions {

		switch instruction {
		case 'L':
			r.turn(-1)
		case 'R':
			r.turn(1)
		case 'F':
			r.step(1)
		case 'B':
			r.step(-1)
		default:

		}

	}

	return r.pos, nil
}

func (r *Rover) turn(turn int) {
	newPos := (r.pos.Dir.GetNumber() + 4 + turn) % 4
	r.pos.Dir = getDirectionByNumber(newPos)
}

func (r *Rover) step(count int) {
	switch r.pos.Dir.GetNumber() {
	case NorthNumber:
		r.pos.Y += count
	case EastNumber:
		r.pos.X += count
	case SouthNumber:
		r.pos.Y -= count
	case WestNumber:
		r.pos.X -= count
	}
}
