package service

import (
	"errors"
	"fmt"
	. "zhasa2.0/branch/entities"
	"zhasa2.0/branch_director/entities"
	"zhasa2.0/branch_director/repo"
	. "zhasa2.0/statistic/entities"
	entities2 "zhasa2.0/user/entities"
)

type BranchDirectorService struct {
	repo repo.BranchDirectorRepository
}

func NewBranchDirectorService(repo repo.BranchDirectorRepository) BranchDirectorService {
	return BranchDirectorService{
		repo: repo,
	}
}

func (bd BranchDirectorService) CreateBranchDirector(userId int32, branchId int32) (entities.BranchDirectorId, error) {
	return bd.repo.CreateBranchDirector(userId, branchId)
}

func (bd BranchDirectorService) SetSmGoal(period Period, smId int32, typeId int32, amount int64) error {
	from, to := period.ConvertToTime()
	return bd.repo.SetSalesManagerGoal(from, to, smId, typeId, amount)
}

func (bd BranchDirectorService) GetBranchDirectorByUserId(userId entities2.UserId) ([]entities.BranchDirector, error) {
	director, err := bd.repo.GetBranchesDirectorByUserId(userId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("director not found by given id userId: ", userId))
	}
	return director, nil
}

func (bd BranchDirectorService) GetBranchDirectorByBranchId(branchId BranchId) (*entities.BranchDirector, error) {
	director, err := bd.repo.GetBranchDirectorByBranchId(branchId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("director not found by given id userId: ", branchId))
	}
	return director, nil
}
