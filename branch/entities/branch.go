package entities

type BranchId int32

type Branch struct {
	BranchId    BranchId
	Title       BranchTitle
	Description BranchDescription
	Key         BranchKey
}

type CreateBranchRequest struct {
	Title       BranchTitle
	Description BranchDescription
	Key         BranchKey
}

type BranchTitle string

func NewBranchTitle(titleText string) BranchTitle {
	return BranchTitle(titleText)
}

type BranchDescription string

func NewBranchDescription(descriptionText string) BranchDescription {
	return BranchDescription(descriptionText)
}

type BranchKey string

func NewBranchKey(keyText string) BranchKey {
	return BranchKey(keyText)
}

type BranchesMap map[BranchId]*Branch
