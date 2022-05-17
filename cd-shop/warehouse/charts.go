package warehouse

//go:generate mockgen -source=charts.go -destination=mock/mock_charts.go -package=mock
type Charts interface {
	Sale(title string, artist string, copies int)
}
