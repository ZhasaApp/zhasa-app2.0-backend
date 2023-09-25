package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type OwnerRepository interface {
}

type DbOwnerRepositroy struct {
	ctx     context.Context
	querier generated.Querier
}

func NewOwnerRepo(ctx context.Context, querier generated.Querier) OwnerRepository {
	return DbOwnerRepositroy{
		ctx:     ctx,
		querier: querier,
	}
}
