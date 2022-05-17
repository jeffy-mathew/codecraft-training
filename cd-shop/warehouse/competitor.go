package warehouse

//go:generate mockgen -source=competitor.go -destination=mock/mock_competitor.go -package=mock

type Competitor interface {
	Price(title string) float64
}
