package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	. "zhasa2.0/db/sqlc"
	. "zhasa2.0/statistic/entities"
	. "zhasa2.0/user/repository"
)

type SaleRepository interface {
	AddOrEdit(saleToCreate AddSaleOrReplaceParams, brandId int32) error
	GetSumByUserIdBrandIdPeriodSaleTypeId(params GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams) (int64, error)
	GetUserBrandMonthlyYearStatistic(year int32, userId int32, brandId int32) ([]MonthlyYearStatistic, error)
	DeleteSale(id int32) error
	GetSalesByBrandIdAndUserId(params GetSalesByBrandIdAndUserIdParams) ([]GetSalesByBrandIdAndUserIdRow, error)
	GetSaleBrandId(saleId int32) (*GetSaleBrandBySaleIdRow, error)
}

type DBSaleRepository struct {
	ctx          context.Context
	store        SaleStore
	brandStore   UserBrandStore
	saleTypeRepo SaleTypeRepository
	brandGoal    UserBrandGoalFunc
	GetUserBrandFunc
}

func (d DBSaleRepository) GetSaleBrandId(saleId int32) (*GetSaleBrandBySaleIdRow, error) {
	saleBrand, err := d.store.GetSaleBrandBySaleId(d.ctx, saleId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &saleBrand, nil
}

func (d DBSaleRepository) GetSalesByBrandIdAndUserId(params GetSalesByBrandIdAndUserIdParams) ([]GetSalesByBrandIdAndUserIdRow, error) {
	rows, err := d.store.GetSalesByBrandIdAndUserId(d.ctx, params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return rows, nil
}

func (d DBSaleRepository) GetUserBrandMonthlyYearStatistic(year int32, userId int32, brandId int32) ([]MonthlyYearStatistic, error) {
	saleTypes, err := d.saleTypeRepo.GetSaleTypes()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := make([]MonthlyYearStatistic, 0)
	userBrand, err := d.brandStore.GetUserBrand(d.ctx, GetUserBrandParams{
		UserID:  userId,
		BrandID: brandId,
	})
	if err != nil {
		fmt.Println(err, "for user id", userId, "and brand id", brandId)
		return nil, errors.New("user brand not found")
	}

	for _, saleType := range *saleTypes {
		for month := 1; month <= 12; month++ {
			period := MonthPeriod{
				MonthNumber: int32(month),
				Year:        year,
			}
			from, to := period.ConvertToTime()
			goal := d.brandGoal(GetUserBrandGoalParams{
				UserBrand:  userBrand,
				SaleTypeID: saleType.Id,
				FromDate:   from,
				FromDate_2: to,
			})

			sum, err := d.store.GetSaleSumByUserIdBrandIdPeriodSaleTypeId(d.ctx, GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams{
				ID:         userId,
				BrandID:    brandId,
				SaleDate:   from,
				SaleDate_2: to,
				SaleTypeID: saleType.Id,
			})
			if err != nil {
				log.Println(err)
			}

			stat := MonthlyYearStatistic{
				SaleType: saleType,
				Month:    int32(month),
				Amount:   sum,
				Goal:     goal,
			}
			result = append(result, stat)
		}
	}
	return result, nil
}

func (d DBSaleRepository) AddOrEdit(saleToCreate AddSaleOrReplaceParams, brandId int32) error {
	_, err := d.store.AddBrandSaleTx(d.ctx, saleToCreate, brandId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (d DBSaleRepository) GetSumByUserIdBrandIdPeriodSaleTypeId(params GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams) (int64, error) {
	amount, err := d.store.GetSaleSumByUserIdBrandIdPeriodSaleTypeId(d.ctx, params)
	if err == sql.ErrNoRows {
		log.Println(err)
		return 0, nil
	}
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return amount, err
}

func (d DBSaleRepository) DeleteSale(id int32) error {
	err := d.store.DeleteSale(d.ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func NewSaleRepo(ctx context.Context, store *DBStore, saleTypeRepo SaleTypeRepository, goalFunc UserBrandGoalFunc) SaleRepository {
	return DBSaleRepository{
		ctx:          ctx,
		store:        store,
		brandStore:   store,
		brandGoal:    goalFunc,
		saleTypeRepo: saleTypeRepo,
	}
}
