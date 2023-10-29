package sm

import "zhasa2.0/user/entities"

type SalesManagerRatingItem struct {
	ID                     int32   `json:"id"`
	Avatar                 *string `json:"avatar"`
	FullName               string  `json:"full_name"`
	Branch                 string  `json:"branch"`
	GoalAchievementPercent float64 `json:"goal_achievement_percent"`
}

type SalesManagerRatingItemsResponse struct {
	Result  []SalesManagerRatingItem `json:"result"`
	Count   int32                    `json:"count"`
	HasNext bool                     `json:"has_next"`
}

func SalesManagerRatingItemFrom(ratedUser entities.RatedUser) SalesManagerRatingItem {
	return SalesManagerRatingItem{
		ID:                     ratedUser.User.Id,
		Avatar:                 ratedUser.AvatarPointer(),
		FullName:               ratedUser.GetFullName(),
		Branch:                 ratedUser.BranchInfo.Title,
		GoalAchievementPercent: ratedUser.Ratio,
	}
}

func SalesManagerRatingItemsFrom(ratedUsers []entities.RatedUser) []SalesManagerRatingItem {
	res := make([]SalesManagerRatingItem, 0)
	for _, ratedUser := range ratedUsers {
		res = append(res, SalesManagerRatingItemFrom(ratedUser))
	}
	return res
}
