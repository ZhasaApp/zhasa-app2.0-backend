package repo

import (
	"context"
	. "zhasa2.0/branch/entities"
	"zhasa2.0/branch_director/entities"
	generated "zhasa2.0/db/sqlc"
	entities2 "zhasa2.0/user/entities"
)

type BranchDirectorRepository interface {
	CreateBranchDirector(userId entities2.UserId, branchId BranchId) (entities.BranchDirectorId, error)
	CreateSalesManagerGoal(goal entities.SalesManagerGoal) error
	GetBranchDirectorByUserId(userId entities2.UserId) (*entities.BranchDirector, error)
}

func NewBranchDirectorRepository(ctx context.Context, querier generated.Querier) BranchDirectorRepository {
	return DbBranchDirectorRepository{
		ctx:     ctx,
		querier: querier,
	}
}

type DbBranchDirectorRepository struct {
	ctx     context.Context
	querier generated.Querier
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

func (bdr DbBranchDirectorRepository) CreateSalesManagerGoal(goal entities.SalesManagerGoal) error {
	params := generated.CreateSalesManagerGoalByTypeParams{
		FromDate:       goal.FromDate,
		ToDate:         goal.ToDate,
		Amount:         int64(goal.Amount),
		SalesManagerID: int32(goal.SalesManagerId),
	}
	return bdr.querier.CreateSalesManagerGoalByType(bdr.ctx, params)
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
