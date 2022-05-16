package warehouse

import (
	"errors"
	"sync"
)

var (
	ErrOutOfStock = errors.New("insufficient stock")
)

type CD struct {
	Title  string
	Artist string
}

type Warehouse struct {
	mu    sync.Mutex
	Stock map[string]int
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
}

func (warehouse *Warehouse) GetStock(title string) int {
	count := warehouse.Stock[title]

	return count
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
