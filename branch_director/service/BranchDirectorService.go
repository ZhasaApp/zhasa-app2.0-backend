package service

import (
	"errors"
	entities3 "zhasa2.0/branch/entities"
	"zhasa2.0/branch_director/entities"
	"zhasa2.0/branch_director/repo"
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

func (bd BranchDirectorService) CreateBranchDirector(userId entities2.UserId, branchId entities3.BranchId) (entities.BranchDirectorId, error) {
	return bd.repo.CreateBranchDirector(userId, branchId)
}

func (bd BranchDirectorService) CreateSalesManagerGoal(goal entities.SalesManagerGoal) error {
	err := bd.repo.CreateSalesManagerGoal(goal)
	if err != nil {
		return errors.New("db error")
	}
	return nil
}

func (bd BranchDirectorService) GetBranchDirectorByUserId(userId entities2.UserId) (*entities.BranchDirector, error) {
	director, err := bd.repo.GetBranchDirectorByUserId(userId)
	if err != nil {
		return nil, errors.New("director not found by given id")
	}
	return director, nil
}
