package warehouse

//go:generate mockgen -source=rank.go -destination=mock/mock_rank.go -package=mock

type Rank interface {
	Get(title string) int
}
