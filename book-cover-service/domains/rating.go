package domains

type Rating struct {
	id       int64
	likes    int
	comments []string

	coverId int64
}
