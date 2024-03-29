package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	. "zhasa2.0/db/sqlc"
	. "zhasa2.0/statistic"
	. "zhasa2.0/statistic/entities"
	. "zhasa2.0/user/repository"
)

type SaleRepository interface {
	GetUserBrandMonthlyYearStatistic(year int32, userId int32, brandId int32) ([]MonthlyYearStatistic, error)
	DeleteSale(id int32) error
	GetSaleBrandId(saleId int32) (*GetSaleBrandBySaleIdRow, error)
}

type DBSaleRepository struct {
	ctx          context.Context
	store        SaleStore
	brandStore   UserBrandStore
	saleTypeRepo SaleTypeRepository
	brandGoal    UserBrandGoalFunc
	userSaleSum  GetSaleSumByUserBrandTypePeriodFunc
}

func (d DBSaleRepository) GetSaleBrandId(saleId int32) (*GetSaleBrandBySaleIdRow, error) {
	saleBrand, err := d.store.GetSaleBrandBySaleId(d.ctx, saleId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &saleBrand, nil
}

func (d DBSaleRepository) GetUserBrandMonthlyYearStatistic(year int32, userId int32, brandId int32) ([]MonthlyYearStatistic, error) {
	saleTypes, err := d.saleTypeRepo.GetSaleTypes()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := make([]MonthlyYearStatistic, 0)

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
			goal := d.brandGoal(
				userId, brandId, saleType.Id, period,
			)

			sum, err := d.userSaleSum(userId, brandId, saleType.Id, period)

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

func (d DBSaleRepository) DeleteSale(id int32) error {
	err := d.store.DeleteSale(d.ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func NewSaleRepo(ctx context.Context, store *DBStore, saleTypeRepo SaleTypeRepository, goalFunc UserBrandGoalFunc, userSaleSumFunc GetSaleSumByUserBrandTypePeriodFunc) SaleRepository {
	return DBSaleRepository{
		ctx:          ctx,
		store:        store,
		brandStore:   store,
		saleTypeRepo: saleTypeRepo,
		brandGoal:    goalFunc,
		userSaleSum:  userSaleSumFunc,
	}
}
