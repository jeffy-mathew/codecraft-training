package rover

import "errors"

var (
	ErrNotDropped = errors.New("rover not dropped")
)

type Direction interface {
	GetNumber() int
}

type North struct{}

func (n North) GetNumber() int {
	return 0
}

type East struct{}

func (e East) GetNumber() int {
	return 1
}

type South struct{}

func (s South) GetNumber() int {
	return 2
}

type West struct{}

func (w West) GetNumber() int {
	return 3
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

	//for _, instruction := range instructions {
	//
	//	switch instruction {
	//	case 'L':
	//		if r.pos.Dir == North {
	//
	//		} else if r.pos.Dir == West
	//	case 'R':
	//	case 'F':
	//	case 'B':
	//	default:
	//
	//	}
	//
	//}

	return r.pos, nil
}
