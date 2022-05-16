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
}

type Warehouse struct {
	mu        sync.Mutex
	Stock     map[string]int
	ArtistCDs map[string][]CD
	TitleCDs  map[string]CD
}

func (warehouse *Warehouse) Add(cd CD, copies int) {
	warehouse.mu.Lock()
	defer warehouse.mu.Unlock()

	if warehouse.Stock == nil {
		warehouse.Stock = make(map[string]int)
	}

	if cdCount, ok := warehouse.Stock[cd.Title]; ok {
		warehouse.Stock[cd.Title] = cdCount + copies
	} else {
		warehouse.Stock[cd.Title] = copies
	}

	if warehouse.ArtistCDs == nil {
		warehouse.ArtistCDs = make(map[string][]CD)
	}

	if artistCDs, ok := warehouse.ArtistCDs[cd.Artist]; ok {
		newArtistCDs := append(artistCDs, cd)
		warehouse.ArtistCDs[cd.Artist] = newArtistCDs
	} else {
		warehouse.ArtistCDs[cd.Artist] = []CD{cd}
	}

	if warehouse.TitleCDs == nil {
		warehouse.TitleCDs = make(map[string]CD)
	}

	if _, ok := warehouse.TitleCDs[cd.Artist]; ok {
		warehouse.TitleCDs[cd.Title] = cd
	} else {
		warehouse.TitleCDs[cd.Title] = cd
	}
}

func (warehouse *Warehouse) GetStock(title string) int {
	warehouse.mu.Lock()
	defer warehouse.mu.Unlock()
	count := warehouse.Stock[title]
	return count
}

func (warehouse *Warehouse) SearchByTitle(title string) (CD, error) {
	warehouse.mu.Lock()
	defer warehouse.mu.Unlock()
	if cd, ok := warehouse.TitleCDs[title]; ok {
		return cd, nil
	}
	return CD{}, ErrCDNotFound
}

func (warehouse *Warehouse) SearchByArtist(artist string) ([]CD, error) {
	warehouse.mu.Lock()
	defer warehouse.mu.Unlock()
	if cds, ok := warehouse.ArtistCDs[artist]; ok {
		return cds, nil
	}
	return nil, ErrCDNotFound
}

func (warehouse *Warehouse) RemoveCDs(title string, copies int) error {
	warehouse.mu.Lock()
	defer warehouse.mu.Unlock()

	if cdCount, ok := warehouse.Stock[title]; ok {
		if cdCount < copies {
			return ErrOutOfStock
		}

		warehouse.Stock[title] = cdCount - copies
	} else {
		return ErrOutOfStock
	}

	return nil
}

func (warehouse *Warehouse) Sell(processor PaymentProcessor, title string, copies int) error {
	cd, err := warehouse.SearchByTitle(title)
	if err != nil {
		return err
	}

	warehouse.mu.Lock()
	cdCount, ok := warehouse.Stock[title]
	warehouse.mu.Unlock()
	if !ok || cdCount < copies {
		return ErrOutOfStock
	}

	totalAmount := float64(copies) * cd.Price

	if err = processor.Pay(totalAmount); err != nil {
		return ErrPaymentFailed
	}

	err = warehouse.RemoveCDs(title, copies)

	return err
}
