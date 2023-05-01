package generated

import "database/sql"

type Store interface {
	Querier
}

type DBStore struct {
	*Queries
}

func NewStore(db *sql.DB) *DBStore {
	return &DBStore{
		New(db),
	}
}
