package position

func GetLatest() *Position {
	return &Position{}
}

func (*Position) Save() {
	/// save the position in db
}
