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
	mu        sync.Mutex
	Stock     map[string]int
	ArtistCDs map[string]int
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
		warehouse.ArtistCDs = make(map[string]int)
	}

	if artistCount, ok := warehouse.ArtistCDs[cd.Artist]; ok {
		warehouse.ArtistCDs[cd.Artist] = artistCount + copies
	} else {
		warehouse.ArtistCDs[cd.Artist] = copies
	}
}

func (warehouse *Warehouse) GetStock(title string) int {
	count := warehouse.Stock[title]

	return count
}

func (warehouse *Warehouse) SearchByTitle(title string) int {
	return warehouse.GetStock(title)
}

func (warehouse *Warehouse) SearchByArtist(artist string) int {
	count := warehouse.ArtistCDs[artist]

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
