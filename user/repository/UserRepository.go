package repository

import (
	db_generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

const (
	connStr = "user=postgres password=yourpassword dbname=yourdbname sslmode=disable"
)

// Email and Password types already provided

type UserRepository interface {
	GetUserByPhone(email entities.Phone) (*entities.User, error)
}

type PostgresUserRepository struct {
	querier db_generated.Querier
}
