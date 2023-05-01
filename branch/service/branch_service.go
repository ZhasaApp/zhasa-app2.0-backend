package service

import (
	"zhasa2.0/branch/entities"
	"zhasa2.0/branch/repository"
)

type BranchService interface {
	CreateBranch(request entities.CreateBranchRequest) error
}

type DBBranchService struct {
	repo repository.BranchRepository
}

func (ds DBBranchService) CreateBranch(request entities.CreateBranchRequest) error {
	return ds.repo.CreateBranch(request)
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return DBBranchService{
		repo: repo,
	}
}
