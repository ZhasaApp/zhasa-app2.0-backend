package repo

import (
	"context"
	"time"
	. "zhasa2.0/branch/entities"
	. "zhasa2.0/branch_director/entities"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/user/entities"
)

type BranchDirectorRepository interface {
	CreateBranchDirector(userId int32, branchId int32) (BranchDirectorId, error)
	GetBranchesDirectorByUserId(userId UserId) ([]BranchDirector, error)
	SetSalesManagerGoal(from, to time.Time, smId int32, saleTypeId int32, amount int64) error
	GetBranchDirectorByBranchId(branch BranchId) (*BranchDirector, error)
	SetBranchGoal(from, to time.Time, branchId int32, saleTypeId int32, amount int64) error
}

func NewBranchDirectorRepository(ctx context.Context, querier generated.Querier) BranchDirectorRepository {
	return DbBranchDirectorRepository{
		ctx:     ctx,
		querier: querier,
	}
}

func (bdr DbBranchDirectorRepository) SetSalesManagerGoal(from, to time.Time, smId int32, saleTypeId int32, amount int64) error {
	params := generated.SetSmGoalBySaleTypeParams{
		FromDate:       from,
		ToDate:         to,
		Amount:         amount,
		SalesManagerID: smId,
		TypeID:         saleTypeId,
	}

	return bdr.querier.SetSmGoalBySaleType(bdr.ctx, params)
}

type DbBranchDirectorRepository struct {
	ctx     context.Context
	querier generated.Querier
}

func (bdr DbBranchDirectorRepository) SetBranchGoal(from, to time.Time, branchId int32, saleTypeId int32, amount int64) error {
	params := generated.SetBranchGoalBySaleTypeParams{
		FromDate: from,
		ToDate:   to,
		Amount:   amount,
		BranchID: branchId,
		TypeID:   saleTypeId,
	}

	return bdr.querier.SetBranchGoalBySaleType(bdr.ctx, params)
}

func (bdr DbBranchDirectorRepository) GetBranchDirectorByBranchId(branch BranchId) (*BranchDirector, error) {
	data, err := bdr.querier.GetBranchDirectorByBranchId(bdr.ctx, int32(branch))
	if err != nil {
		return nil, err
	}

	director := BranchDirector{
		User: User{
			Id:        data.UserID,
			Phone:     Phone(data.Phone),
			Avatar:    data.AvatarUrl,
			FirstName: Name(data.FirstName),
			LastName:  Name(data.LastName),
		},
		BranchDirectorId: BranchDirectorId(data.BranchDirectorID),
		Branch: Branch{
			BranchId:    BranchId(data.BranchID),
			Title:       BranchTitle(data.BranchTitle),
			Description: "",
			Key:         "",
		},
	}
	return &director, nil
}

func (bdr DbBranchDirectorRepository) CreateBranchDirector(userId int32, branchId int32) (BranchDirectorId, error) {
	params := generated.CreateBranchDirectorParams{
		UserID:   userId,
		BranchID: branchId,
	}
	id, err := bdr.querier.CreateBranchDirector(bdr.ctx, params)
	if err != nil {
		return -1, err
	}
	return BranchDirectorId(id), nil
}

func (bdr DbBranchDirectorRepository) GetBranchesDirectorByUserId(userId UserId) ([]BranchDirector, error) {
	data, err := bdr.querier.GetBranchDirectorByUserId(bdr.ctx, int32(userId))
	if err != nil {
		return nil, err
	}

	directors := make([]BranchDirector, 0)
	for _, row := range data {
		director := BranchDirector{
			User: User{
				Id:        row.UserID,
				Phone:     Phone(row.Phone),
				Avatar:    row.AvatarUrl,
				FirstName: Name(row.FirstName),
				LastName:  Name(row.LastName),
			},
			BranchDirectorId: BranchDirectorId(row.BranchDirectorID),
			Branch: Branch{
				BranchId:    BranchId(row.BranchID),
				Title:       BranchTitle(row.BranchTitle),
				Description: "",
				Key:         "",
			},
		}
		directors = append(directors, director)

	}
	return directors, nil
}
