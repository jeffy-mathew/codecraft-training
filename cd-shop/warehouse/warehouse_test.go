package warehouse

import (
	mock_warehouse "codecraft/cd-shop/warehouse/mock"
	"github.com/golang/mock/gomock"
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
	darkSide := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd", Price: 30.0}
	warehouse.Add(darkSide, 10)

	brainDamage := CD{Title: "Brain damage", Artist: "Pink Floyd", Price: 45.0}
	warehouse.Add(brainDamage, 20)

	breathe := CD{Title: "Breathe", Artist: "Pink Floyd", Price: 25.0}
	warehouse.Add(breathe, 30)

	t.Run("search by title", func(t *testing.T) {
		cd, err := warehouse.SearchByTitle(breathe.Title)
		assert.NoError(t, err)
		assert.Equal(t, breathe, cd)
	})

	t.Run("no CD found for title", func(t *testing.T) {
		cd, err := warehouse.SearchByTitle("Country Song")
		assert.ErrorIs(t, err, ErrCDNotFound)
		assert.Empty(t, cd)
	})

	t.Run("no CDs found for artist", func(t *testing.T) {
		cd, err := warehouse.SearchByArtist("avicii")
		assert.ErrorIs(t, err, ErrCDNotFound)
		assert.Empty(t, cd)
	})

	t.Run("search by artist", func(t *testing.T) {
		artistCDs, err := warehouse.SearchByArtist(brainDamage.Artist)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []CD{brainDamage, darkSide, breathe}, artistCDs)
	})

}

func Test_WarehouseSell(t *testing.T) {
	warehouse := Warehouse{}
	darkSide := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd", Price: 30.0}
	warehouse.Add(darkSide, 10)

	t.Run("fail when title cd not found", func(t *testing.T) {
		err := warehouse.Sell(CreditCard{}, "Closer", 10)
		assert.ErrorIs(t, err, ErrCDNotFound)
	})

	t.Run("fail when out of stock", func(t *testing.T) {
		err := warehouse.Sell(CreditCard{}, darkSide.Title, 100)
		assert.ErrorIs(t, err, ErrOutOfStock)
	})

	t.Run("accept payment and reduce stock", func(t *testing.T) {
		err := warehouse.Sell(CreditCard{}, darkSide.Title, 10)
		assert.NoError(t, err)

		totalCDsLeft := warehouse.GetStock(darkSide.Title)
		assert.Equal(t, 0, totalCDsLeft)
	})

	t.Run("do not sell when payment fails", func(t *testing.T) {
		warehouse.Add(darkSide, 10)
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()

		errPaymentProcessor := mock_warehouse.NewMockPaymentProcessor(ctrl)
		errPaymentProcessor.EXPECT().Pay(300.0).Return(ErrPaymentFailed)

		err := warehouse.Sell(errPaymentProcessor, darkSide.Title, 10)
		assert.ErrorIs(t, ErrPaymentFailed, err)
	})
}
