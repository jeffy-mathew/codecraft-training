package warehouse

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WarehouseAddCDs(t *testing.T) {
	t.Run("add a copies of a CD", func(t *testing.T) {
		warehouse := Warehouse{}

		cd := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd"}
		warehouse.Add(cd, 1)
		totalCDs := warehouse.GetStock(cd.Title)

		assert.Equal(t, 1, totalCDs)
	})

	t.Run("add 10 copies of a CD", func(t *testing.T) {
		warehouse := Warehouse{}

		cd := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd"}
		warehouse.Add(cd, 10)
		totalCDs := warehouse.GetStock(cd.Title)

		assert.Equal(t, 10, totalCDs)
	})
}

func Test_WarehouseRemoveCDs(t *testing.T) {
	t.Run("remove a copy of a CD", func(t *testing.T) {
		warehouse := Warehouse{}

		cd := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd"}
		warehouse.Add(cd, 10)
		err := warehouse.RemoveCDs(cd.Title, 1)
		assert.NoError(t, err)

	})

	t.Run("remove a copy of CD whose stock is empty", func(t *testing.T) {
		warehouse := Warehouse{}

		cd := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd"}

		err := warehouse.RemoveCDs(cd.Title, 1)

		assert.ErrorIs(t, err, ErrOutOfStock)
	})

	t.Run("remove copies more than in stock ", func(t *testing.T) {
		warehouse := Warehouse{}

		cd := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd"}
		warehouse.Add(cd, 10)

		err := warehouse.RemoveCDs(cd.Title, 11)
		assert.ErrorIs(t, err, ErrOutOfStock)
	})
}

func Test_WarehouseSearchCD(t *testing.T) {
	warehouse := Warehouse{}
	darkSide := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd"}
	warehouse.Add(darkSide, 10)

	brainDamage := CD{Title: "Brain damage", Artist: "Pink Floyd"}
	warehouse.Add(brainDamage, 20)

	breathe := CD{Title: "Breathe", Artist: "Pink Floyd"}
	warehouse.Add(breathe, 30)

	t.Run("search by title", func(t *testing.T) {

		copies := warehouse.SearchByTitle(breathe.Title)
		assert.Equal(t, 30, copies)
	})

	t.Run("search by artist", func(t *testing.T) {
		copies := warehouse.SearchByArtist(brainDamage.Artist)
		assert.Equal(t, 60, copies)
	})

}
