package entities

type BranchId int32

type Branch struct {
	BranchId        int32
	Title           string
	Description     string
	Key             string
	GoalAchievement float32
}

type CreateBranchRequest struct {
	Title       BranchTitle
	Description BranchDescription
	Key         BranchKey
}

type BranchTitle string

type BranchDescription string

type BranchKey string

type BranchesMap map[BranchId]*Branch

type BranchDescriptionInfo struct {
	BranchId    int32
	Title       string
	Description string
}
