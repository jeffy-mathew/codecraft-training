package rover

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDrop(t *testing.T) {
	t.Run("drop to position with valid direction", func(t *testing.T) {
		rover := Rover{}
		position, err := rover.Drop(2, 4, North{})

		expectedPosition := &Position{2, 4, North{}}
		assert.NoError(t, err)
		assert.Equal(t, expectedPosition, position)
	})

}

func TestMove(t *testing.T) {

	t.Run("before drop", func(t *testing.T) {
		rover := Rover{}
		position, err := rover.Move("")

		assert.ErrorIs(t, err, ErrNotDropped)
		assert.Nil(t, position)
	})

	t.Run("after drop", func(t *testing.T) {
		rover := Rover{}
		_, _ = rover.Drop(2, 5, North{})

		position, err := rover.Move("")
		expectedPosition := &Position{2, 5, North{}}

		assert.NoError(t, err)
		assert.Equal(t, expectedPosition, position)
	})

	t.Run("turn left", func(t *testing.T) {
		rover := Rover{}
		_, _ = rover.Drop(2, 5, North{})

		position, err := rover.Move("L")
		expectedPosition := &Position{2, 5, West{}}

		assert.NoError(t, err)
		assert.Equal(t, expectedPosition, position)
	})

	t.Run("turn left twice", func(t *testing.T) {
		rover := Rover{}
		_, _ = rover.Drop(2, 5, North{})

		position, err := rover.Move("LL")
		expectedPosition := &Position{2, 5, South{}}

		assert.NoError(t, err)
		assert.Equal(t, expectedPosition, position)
	})

	t.Run("turn right", func(t *testing.T) {
		rover := Rover{}
		_, _ = rover.Drop(2, 5, South{})

		position, err := rover.Move("R")
		expectedPosition := &Position{2, 5, West{}}

		assert.NoError(t, err)
		assert.Equal(t, expectedPosition, position)
	})

	t.Run("move forward from south", func(t *testing.T) {
		rover := Rover{}
		_, _ = rover.Drop(2, 5, South{})

		position, err := rover.Move("F")
		expectedPosition := &Position{2, 4, South{}}

		assert.NoError(t, err)
		assert.Equal(t, expectedPosition, position)
	})

	t.Run("move forward from north", func(t *testing.T) {
		rover := Rover{}
		_, _ = rover.Drop(2, 5, North{})

		position, err := rover.Move("F")
		expectedPosition := &Position{2, 6, North{}}

		assert.NoError(t, err)
		assert.Equal(t, expectedPosition, position)
	})

	t.Run("move forward from west", func(t *testing.T) {
		rover := Rover{}
		_, _ = rover.Drop(2, 5, West{})

		position, err := rover.Move("F")
		expectedPosition := &Position{1, 5, West{}}

		assert.NoError(t, err)
		assert.Equal(t, expectedPosition, position)
	})

	t.Run("move forward from east", func(t *testing.T) {
		rover := Rover{}
		_, _ = rover.Drop(2, 5, East{})

		position, err := rover.Move("F")
		expectedPosition := &Position{3, 5, East{}}

		assert.NoError(t, err)
		assert.Equal(t, expectedPosition, position)
	})

}

func TestTurn(t *testing.T) {
	rover := Rover{}
	_, _ = rover.Drop(2, 5, North{})

	t.Run("left", func(t *testing.T) {

		rover.turn(-1)
		assert.Equal(t, West{}, rover.pos.Dir)
	})

}
