package bt

import "algo/src/core/db"

type Repo struct {
	db db.Database
}

func NewRepo(
	factory *db.DBFactory,
) *Repo {
	return &Repo{
		db: factory.Get("clickhouse"),
	}
}

func (this *Repo) Create() *Run {
	// TODO: do this
	this.db.Insert("INSERT INTO bt_run VALUES ()")

	res, err := this.db.Select("SELECT FROM bt_run WHERE id = '{id}'")

	if err != nil {
		panic("")
	}

	v, _ := res[0].(Run)

	return &v
}

func (*Repo) Get() {

}

func (*Repo) GetLatest() {

}

func (*Repo) Update() {

}

func (*Repo) Delete() {

}
