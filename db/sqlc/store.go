package generated

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
}

type DBStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) *DBStore {
	return &DBStore{
		db,
		New(db),
	}
}

func (store *DBStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *DBStore) AddBrandSaleTx(ctx context.Context, params AddSaleOrReplaceParams, brandId int32) (*Sale, error) {
	var sale *Sale

	err := store.execTx(ctx, func(queries *Queries) error {
		res, err := queries.AddSaleOrReplace(ctx, params)
		if err != nil {
			return err
		}
		_, err = queries.AddSaleToBrand(ctx, AddSaleToBrandParams{
			SaleID:  res.ID,
			BrandID: brandId,
		})
		if err != nil {
			fmt.Println(err, " ", params, " ", brandId)
			return err
		}
		sale = &res
		return nil
	})

	return sale, err
}
