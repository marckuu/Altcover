package domains

type Comment struct {
	ID   int64
	Text string

	CoverId int64
	UserId  int64
}
