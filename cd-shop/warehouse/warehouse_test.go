package warehouse

import (
	"codecraft/cd-shop/warehouse/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WarehouseAddCDs(t *testing.T) {
	t.Run("add a copies of a CD", func(t *testing.T) {
		warehouse := Warehouse{}

		cd := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd", stock: 1}
		warehouse.Add(&cd)
		totalCDs := cd.GetStock()

		assert.Equal(t, 1, totalCDs)
	})

	t.Run("add 10 copies of a CD", func(t *testing.T) {
		warehouse := Warehouse{}

		cd := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd", stock: 10}
		warehouse.Add(&cd)
		totalCDs := cd.GetStock()

		assert.Equal(t, 10, totalCDs)
	})
}

func Test_WarehouseRemoveCDs(t *testing.T) {
	t.Run("remove a copy of a CD", func(t *testing.T) {
		warehouse := Warehouse{}

		cd := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd", stock: 10}
		warehouse.Add(&cd)
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

		cd := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd", stock: 10}
		warehouse.Add(&cd)

		err := warehouse.RemoveCDs(cd.Title, 11)
		assert.ErrorIs(t, err, ErrOutOfStock)
	})
}

func Test_WarehouseSearchCD(t *testing.T) {
	warehouse := Warehouse{}
	darkSide := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd", Price: 30.0, stock: 10}
	warehouse.Add(&darkSide)

	brainDamage := CD{Title: "Brain damage", Artist: "Pink Floyd", Price: 45.0, stock: 20}
	warehouse.Add(&brainDamage)

	breathe := CD{Title: "Breathe", Artist: "Pink Floyd", Price: 25.0, stock: 30}
	warehouse.Add(&breathe)

	t.Run("search by title", func(t *testing.T) {
		cd, err := warehouse.SearchByTitle(breathe.Title)
		assert.NoError(t, err)
		assert.Equal(t, breathe, *cd)
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
		artistCDs, err := warehouse.SearchByArtist("Pink Floyd")
		assert.NoError(t, err)
		assert.ElementsMatch(t, []*CD{&brainDamage, &darkSide, &breathe}, artistCDs)
	})

}

func Test_WarehouseSell(t *testing.T) {
	warehouse := Warehouse{}
	darkSide := CD{Title: "The Dark Side of the Moon", Artist: "Pink Floyd", Price: 30.0, stock: 10}
	warehouse.Add(&darkSide)

	t.Run("fail when title cd not found", func(t *testing.T) {
		err := warehouse.Sell(CreditCard{}, &CD{Title: "Closer"}, 10)
		assert.ErrorIs(t, err, ErrCDNotFound)
	})

	t.Run("fail when out of stock", func(t *testing.T) {
		err := warehouse.Sell(CreditCard{}, &CD{Title: "The Dark Side of the Moon"}, 100)
		assert.ErrorIs(t, err, ErrOutOfStock)
	})

	t.Run("accept payment and reduce stock", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()

		paymentProcessor := mock.NewMockPaymentProcessor(ctrl)
		paymentProcessor.EXPECT().Pay(300.0).Return(nil)

		err := warehouse.Sell(paymentProcessor, &CD{Title: "The Dark Side of the Moon"}, 10)
		assert.NoError(t, err)

		totalCDsLeft := darkSide.GetStock()
		assert.Equal(t, 0, totalCDsLeft)
	})

	t.Run("notify chart of sales", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		chartMock := mock.NewMockCharts(ctrl)

		newDarkSide := darkSide
		warehouse := Warehouse{
			chart: chartMock,
		}

		warehouse.Add(&newDarkSide)
		copies := 2

		chartMock.EXPECT().Sale(newDarkSide.Title, newDarkSide.Artist, copies)

		err := warehouse.Sell(CreditCard{}, &newDarkSide, copies)
		assert.NoError(t, err)
	})

	t.Run("do not sell when payment fails", func(t *testing.T) {
		newDarkSide := darkSide
		newDarkSide.stock = 10
		warehouse.Add(&newDarkSide)
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()

		errPaymentProcessor := mock.NewMockPaymentProcessor(ctrl)
		errPaymentProcessor.EXPECT().Pay(300.0).Return(ErrPaymentFailed)

		err := warehouse.Sell(errPaymentProcessor, &CD{Title: "The Dark Side of the Moon"}, 10)
		assert.ErrorIs(t, ErrPaymentFailed, err)
	})
}
