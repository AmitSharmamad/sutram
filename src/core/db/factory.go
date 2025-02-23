package db

type DBFactory struct {
	dbs map[string]Database
}

func NewFactoryProv(
	clickhouse *Clickhouse,
	postgres *Postgres,
) *DBFactory {
	factory := &DBFactory{
		dbs: make(map[string]Database),
	}
	factory.dbs[string(CLICKHOUSE)] = clickhouse
	factory.dbs[string(POSTGRES)] = postgres
	return factory
}

func (factory *DBFactory) Register(key string, db Database) {
	factory.dbs[key] = db
}

func (factory *DBFactory) Get(key string) Database {
	return factory.dbs[key]
}
