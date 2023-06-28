// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package generated

import (
	"time"
)

type Branch struct {
	ID          int32     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	BranchKey   string    `json:"branch_key"`
	CreatedAt   time.Time `json:"created_at"`
}

type BranchDirector struct {
	ID       int32 `json:"id"`
	UserID   int32 `json:"user_id"`
	BranchID int32 `json:"branch_id"`
}

type BranchDirectorsView struct {
	UserID           int32  `json:"user_id"`
	Phone            string `json:"phone"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	AvatarUrl        string `json:"avatar_url"`
	BranchDirectorID int32  `json:"branch_director_id"`
	BranchID         int32  `json:"branch_id"`
	BranchTitle      string `json:"branch_title"`
}

type BranchGoal struct {
	ID       int32     `json:"id"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	Amount   int64     `json:"amount"`
	BranchID int32     `json:"branch_id"`
}

type Sale struct {
	ID             int32     `json:"id"`
	SalesManagerID int32     `json:"sales_manager_id"`
	SaleDate       time.Time `json:"sale_date"`
	Amount         int64     `json:"amount"`
	SaleTypeID     int32     `json:"sale_type_id"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
}

type SaleType struct {
	ID          int32     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Color       string    `json:"color"`
	Gravity     int32     `json:"gravity"`
}

type SalesManager struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	BranchID  int32     `json:"branch_id"`
	CreatedAt time.Time `json:"created_at"`
}

type SalesManagerGoalsByType struct {
	ID             int32     `json:"id"`
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
	Amount         int64     `json:"amount"`
	SalesManagerID int32     `json:"sales_manager_id"`
	TypeID         int32     `json:"type_id"`
}

type SalesManagerGoalsRatioByPeriod struct {
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
	Ratio          float64   `json:"ratio"`
	SalesManagerID int32     `json:"sales_manager_id"`
}

type SalesManagersView struct {
	UserID         int32  `json:"user_id"`
	Phone          string `json:"phone"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	AvatarUrl      string `json:"avatar_url"`
	SalesManagerID int32  `json:"sales_manager_id"`
	BranchID       int32  `json:"branch_id"`
	BranchTitle    string `json:"branch_title"`
}

type User struct {
	ID        int32     `json:"id"`
	Phone     string    `json:"phone"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}

type UserAvatarView struct {
	ID        int32  `json:"id"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarUrl string `json:"avatar_url"`
}

type UsersAvatar struct {
	UserID    int32  `json:"user_id"`
	AvatarUrl string `json:"avatar_url"`
}

type UsersCode struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Code      int32     `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}
