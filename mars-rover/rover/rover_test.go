package rover

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const invalidDirection Direction = "H"

func TestDrop(t *testing.T) {
	t.Run("drop to unset position", func(t *testing.T) {
		rover := Rover{}
		position, err := rover.Drop(0, 0, "")

		assert.ErrorIs(t, err, ErrInvalidDirection)
		assert.Nil(t, position)
	})

	t.Run("drop to position with invalid direction", func(t *testing.T) {
		rover := Rover{}

		position, err := rover.Drop(2, 4, invalidDirection)

		assert.ErrorIs(t, err, ErrInvalidDirection)
		assert.Nil(t, position)
	})

	t.Run("drop to position with valid direction", func(t *testing.T) {
		rover := Rover{}
		position, err := rover.Drop(2, 4, North)

		expectedPosition := &Position{2, 4, North}
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
		_, _ = rover.Drop(2, 5, North)

		position, err := rover.Move("")
		expectedPosition := &Position{2, 5, North}

		assert.NoError(t, err)
		assert.Equal(t, expectedPosition, position)
	})

	t.Run("move left", func(t *testing.T) {
		//rover := Rover{}
		//_, _ = rover.Drop(2, 5, North)
		//
		//assert.NoError(t, err)
		//assert.Equal(t, expectedPosition, position)
	})
}
