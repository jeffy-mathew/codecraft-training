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
		err := warehouse.Sell(CreditCard{}, &CD{Title: "Closer"}, 10, 0)
		assert.ErrorIs(t, err, ErrCDNotFound)
	})

	t.Run("fail when out of stock", func(t *testing.T) {
		err := warehouse.Sell(CreditCard{}, &CD{Title: "The Dark Side of the Moon"}, 100, 0)
		assert.ErrorIs(t, err, ErrOutOfStock)
	})

	t.Run("accept payment and reduce stock", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()

		paymentProcessor := mock.NewMockPaymentProcessor(ctrl)
		paymentProcessor.EXPECT().Pay(300.0).Return(nil)

		mockRank := mock.NewMockRank(ctrl)
		mockRank.EXPECT().Get(gomock.Any()).Return(120)

		mockCompetitor := mock.NewMockCompetitor(ctrl)

		mockChart := mock.NewMockCharts(ctrl)
		mockChart.EXPECT().Sale(gomock.Any(), gomock.Any(), gomock.Any())

		warehouse.competitor = mockCompetitor
		warehouse.rank = mockRank
		warehouse.chart = mockChart

		err := warehouse.Sell(paymentProcessor, &CD{Title: "The Dark Side of the Moon"}, 10, 0)
		assert.NoError(t, err)

		totalCDsLeft := darkSide.GetStock()
		assert.Equal(t, 0, totalCDsLeft)
	})

	t.Run("notify chart of sales", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		chartMock := mock.NewMockCharts(ctrl)
		rankMock := mock.NewMockRank(ctrl)
		rankMock.EXPECT().Get(gomock.Any()).Return(120)

		breathe := CD{Title: "Breathe", Artist: "Pink Floyd", Price: 25.0, stock: 30}
		warehouse.Add(&breathe)

		warehouse := Warehouse{
			chart: chartMock,
			rank:  rankMock,
		}

		warehouse.Add(&breathe)
		copies := 2

		chartMock.EXPECT().Sale(breathe.Title, breathe.Artist, copies)

		err := warehouse.Sell(CreditCard{}, &breathe, copies, 0)
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
		mockRank := mock.NewMockRank(ctrl)
		mockRank.EXPECT().Get(gomock.Any()).Return(129)
		warehouse.rank = mockRank

		err := warehouse.Sell(errPaymentProcessor, &CD{Title: "The Dark Side of the Moon"}, 10, 0)
		assert.ErrorIs(t, ErrPaymentFailed, err)
	})

	t.Run("offer 1 pound less than competitors", func(t *testing.T) {

		newDarkSide := darkSide
		newDarkSide.stock = 10
		warehouse.Add(&newDarkSide)

		ctrl := gomock.NewController(t)
		mockRank := mock.NewMockRank(ctrl)
		mockRank.EXPECT().Get(newDarkSide.Title).Return(5)

		mockCompetitor := mock.NewMockCompetitor(ctrl)
		mockCompetitor.EXPECT().Price(newDarkSide.Title).Return(20.0)

		chartMock := mock.NewMockCharts(ctrl)
		chartMock.EXPECT().Sale(gomock.Any(), gomock.Any(), gomock.Any())

		warehouse.chart = chartMock
		warehouse.rank = mockRank
		warehouse.competitor = mockCompetitor

		paymentProcessor := mock.NewMockPaymentProcessor(ctrl)
		paymentProcessor.EXPECT().Pay(38.0).Return(nil)

		err := warehouse.Sell(paymentProcessor, &newDarkSide, 2, 0)
		assert.NoError(t, err)
	})

}
