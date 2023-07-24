package repo

import (
	"context"
	"time"
	. "zhasa2.0/branch/entities"
	"zhasa2.0/branch_director/entities"
	generated "zhasa2.0/db/sqlc"
	entities2 "zhasa2.0/user/entities"
)

type BranchDirectorRepository interface {
	CreateBranchDirector(userId entities2.UserId, branchId BranchId) (entities.BranchDirectorId, error)
	GetBranchDirectorByUserId(userId entities2.UserId) (*entities.BranchDirector, error)
	SetSalesManagerGoal(from, to time.Time, smId int32, saleTypeId int32, amount int64) error
	GetBranchDirectorByBranchId(branch BranchId) (*entities.BranchDirector, error)
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

func (bdr DbBranchDirectorRepository) GetBranchDirectorByBranchId(branch BranchId) (*entities.BranchDirector, error) {
	data, err := bdr.querier.GetBranchDirectorByBranchId(bdr.ctx, int32(branch))
	if err != nil {
		return nil, err
	}
	director := entities.BranchDirector{
		User: entities2.User{
			Id:    data.UserID,
			Phone: entities2.Phone(data.Phone),
			Avatar: entities2.Avatar{
				Url: data.AvatarUrl,
			},
			FirstName: entities2.Name(data.FirstName),
			LastName:  entities2.Name(data.LastName),
		},
		BranchDirectorId: entities.BranchDirectorId(data.BranchDirectorID),
		Branch: Branch{
			BranchId:    BranchId(data.BranchID),
			Title:       BranchTitle(data.BranchTitle),
			Description: "",
			Key:         "",
		},
	}
	return &director, nil
}

func (bdr DbBranchDirectorRepository) CreateBranchDirector(userId entities2.UserId, branchId BranchId) (entities.BranchDirectorId, error) {
	params := generated.CreateBranchDirectorParams{
		UserID:   int32(userId),
		BranchID: int32(branchId),
	}
	id, err := bdr.querier.CreateBranchDirector(bdr.ctx, params)
	if err != nil {
		return -1, err
	}
	return entities.BranchDirectorId(id), nil
}

func (bdr DbBranchDirectorRepository) GetBranchDirectorByUserId(userId entities2.UserId) (*entities.BranchDirector, error) {
	data, err := bdr.querier.GetBranchDirectorByUserId(bdr.ctx, int32(userId))
	if err != nil {
		return nil, err
	}
	director := entities.BranchDirector{
		User: entities2.User{
			Id:    data.UserID,
			Phone: entities2.Phone(data.Phone),
			Avatar: entities2.Avatar{
				Url: data.AvatarUrl,
			},
			FirstName: entities2.Name(data.FirstName),
			LastName:  entities2.Name(data.LastName),
		},
		BranchDirectorId: entities.BranchDirectorId(data.BranchDirectorID),
		Branch: Branch{
			BranchId:    BranchId(data.BranchID),
			Title:       BranchTitle(data.BranchTitle),
			Description: "",
			Key:         "",
		},
	}
	return &director, nil
}
