package entities

import "zhasa2.0/base"

type BranchId int32

type Branch struct {
	BranchId        BranchId
	Title           BranchTitle
	Description     BranchDescription
	Key             BranchKey
	GoalAchievement base.Percent
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

type BranchDescriptionInfo struct {
	BranchId    int32
	Title       string
	Description string
}
