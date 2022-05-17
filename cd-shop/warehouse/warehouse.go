package warehouse

//go:generate mockgen -source=warehouse.go -destination=mock/mock_warehouse.go -package=mock

import (
	"errors"
	"sync"
)

var (
	ErrOutOfStock    = errors.New("insufficient stock")
	ErrPaymentFailed = errors.New("payment failed")
	ErrCDNotFound    = errors.New("CD(s) not found")
)

type PaymentProcessor interface {
	Pay(float64) error
}

func (c CreditCard) Pay(f float64) error {
	return nil
}

type CreditCard struct{}

type CD struct {
	Title  string
	Artist string
	Price  float64
	stock  int
}

type Warehouse struct {
	chart      Charts
	rank       Rank
	competitor Competitor
	mu         sync.Mutex
	cds        []*CD
}

func (warehouse *Warehouse) Add(cd *CD) {
	warehouse.mu.Lock()
	defer warehouse.mu.Unlock()

	for i, existingCD := range warehouse.cds {
		if existingCD.Title == cd.Title {
			warehouse.cds[i].stock += cd.stock
			return
		}
	}

	warehouse.cds = append(warehouse.cds, cd)

}

func (cd *CD) GetStock() int {
	return cd.stock
}

func (warehouse *Warehouse) SearchByTitle(title string) (*CD, error) {
	warehouse.mu.Lock()
	defer warehouse.mu.Unlock()
	for _, cd := range warehouse.cds {
		if cd.Title == title {
			return cd, nil
		}
	}
	return nil, ErrCDNotFound
}

func (warehouse *Warehouse) SearchByArtist(artist string) ([]*CD, error) {
	warehouse.mu.Lock()
	defer warehouse.mu.Unlock()

	var artistCDs []*CD

	for i := 0; i < len(warehouse.cds); i++ {
		if warehouse.cds[i].Artist == artist {
			artistCDs = append(artistCDs, warehouse.cds[i])
		}
	}

	if len(artistCDs) > 0 {
		return artistCDs, nil
	}

	return nil, ErrCDNotFound
}

func (warehouse *Warehouse) RemoveCDs(title string, copies int) error {
	warehouse.mu.Lock()
	defer warehouse.mu.Unlock()

	var selectedCD *CD
	for i := 0; i < len(warehouse.cds); i++ {
		if warehouse.cds[i].Title == title {
			selectedCD = warehouse.cds[i]
		}
	}

	if selectedCD == nil || selectedCD.stock < copies {
		return ErrOutOfStock
	}

	selectedCD.stock -= copies

	return nil
}

func (warehouse *Warehouse) Sell(processor PaymentProcessor, cd *CD, copies int, foo int) error {

	warehouse.mu.Lock()
	var selectedCD *CD
	for i := 0; i < len(warehouse.cds); i++ {
		if warehouse.cds[i].Title == cd.Title {
			selectedCD = warehouse.cds[i]
		}
	}
	warehouse.mu.Unlock()

	if selectedCD == nil {
		return ErrCDNotFound
	}

	if selectedCD.stock < copies {
		return ErrOutOfStock
	}

	offerPrice := selectedCD.Price

	rank := warehouse.rank.Get(selectedCD.Title)
	if rank < 100 {
		competitorPrice := warehouse.competitor.Price(selectedCD.Title)
		if competitorPrice < selectedCD.Price {
			offerPrice = competitorPrice - 1
		}
	}

	totalAmount := float64(copies) * offerPrice

	if err := processor.Pay(totalAmount); err != nil {
		return ErrPaymentFailed
	}

	err := warehouse.RemoveCDs(cd.Title, copies)
	if err != nil {
		return err
	}

	warehouse.chart.Sale(selectedCD.Title, selectedCD.Artist, copies)

	return nil
}
