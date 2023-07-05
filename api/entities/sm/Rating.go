package sm

type SalesManagerRatingItem struct {
	ID                     int32   `json:"id"`
	Avatar                 *string `json:"avatar"`
	FullName               string  `json:"full_name"`
	Branch                 string  `json:"branch"`
	GoalAchievementPercent float64 `json:"goal_achievement_percent"`
}

type SalesManagersResponse struct {
	Result  []SalesManagerRatingItem `json:"result"`
	Count   int32                    `json:"count"`
	HasNext bool                     `json:"has_next"`
}
