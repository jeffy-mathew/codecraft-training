package warehouse

//go:generate mockgen -source=warehouse.go -destination=mock/mock_warehouse.go

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
	mu  sync.Mutex
	cds []*CD
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

func (warehouse *Warehouse) Sell(processor PaymentProcessor, title string, copies int) error {
	cd, err := warehouse.SearchByTitle(title)
	if err != nil {
		return err
	}

	warehouse.mu.Lock()
	cdCount := cd.GetStock()
	warehouse.mu.Unlock()
	if cdCount < copies {
		return ErrOutOfStock
	}

	totalAmount := float64(copies) * cd.Price

	if err = processor.Pay(totalAmount); err != nil {
		return ErrPaymentFailed
	}

	err = warehouse.RemoveCDs(title, copies)

	return err
}
