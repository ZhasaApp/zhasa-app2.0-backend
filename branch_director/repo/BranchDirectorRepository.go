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
	CreateBranchDirector(userId UserId, branchId BranchId) (BranchDirectorId, error)
	GetBranchDirectorByUserId(userId UserId) (*BranchDirector, error)
	SetSalesManagerGoal(from, to time.Time, smId int32, saleTypeId int32, amount int64) error
	GetBranchDirectorByBranchId(branch BranchId) (*BranchDirector, error)
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

func (bdr DbBranchDirectorRepository) GetBranchDirectorByBranchId(branch BranchId) (*BranchDirector, error) {
	data, err := bdr.querier.GetBranchDirectorByBranchId(bdr.ctx, int32(branch))
	if err != nil {
		return nil, err
	}
	director := BranchDirector{
		User: User{
			Id:        data.UserID,
			Phone:     Phone(data.Phone),
			Avatar:    &data.AvatarUrl,
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

func (bdr DbBranchDirectorRepository) CreateBranchDirector(userId UserId, branchId BranchId) (BranchDirectorId, error) {
	params := generated.CreateBranchDirectorParams{
		UserID:   int32(userId),
		BranchID: int32(branchId),
	}
	id, err := bdr.querier.CreateBranchDirector(bdr.ctx, params)
	if err != nil {
		return -1, err
	}
	return BranchDirectorId(id), nil
}

func (bdr DbBranchDirectorRepository) GetBranchDirectorByUserId(userId UserId) (*BranchDirector, error) {
	data, err := bdr.querier.GetBranchDirectorByUserId(bdr.ctx, int32(userId))
	if err != nil {
		return nil, err
	}
	director := BranchDirector{
		User: User{
			Id:        data.UserID,
			Phone:     Phone(data.Phone),
			Avatar:    &data.AvatarUrl,
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
