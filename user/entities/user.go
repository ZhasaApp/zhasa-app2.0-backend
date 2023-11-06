package entities

import (
	"errors"
	"strings"
	"time"
)

type CreateUserRequest struct {
	Phone     Phone
	FirstName Name
	LastName  Name
}

type User struct {
	Id        int32
	Phone     Phone
	Avatar    string
	FirstName string
	LastName  string
	UserRole  UserRole
}

func (u User) AvatarPointer() *string {
	if len(u.Avatar) == 0 {
		return nil
	}
	return &u.Avatar
}

func (u User) GetFullName() string {
	return strings.TrimSpace(u.FirstName) + " " + strings.TrimSpace(u.LastName)
}

type UserAuth struct {
	Code      OtpCode
	UserId    UserId
	CreatedAt time.Time
}

type UserId int32

type Name string

type OtpCode int32

type OtpId int32

func NewName(name string) (*Name, error) {
	// Check that the name is not empty
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	res := Name(name)
	return &res, nil
}

func (n Name) String() string {
	return string(n)
}

type UserRole struct {
	Id  int32
	Key string
}

type RatedUser struct {
	User
	Ratio float64
	BranchInfo
}

type BranchInfo struct {
	Id    int32
	Title string
}

type UserWithBrands struct {
	Id          int32  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	BranchTitle string `json:"branch_title"`
	Brands      string `json:"brands"`
}
