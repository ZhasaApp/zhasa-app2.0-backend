package entities

import "fmt"

type SaleAmount int64

type Percent float64

func NewPercent(achieved SaleAmount, goal SaleAmount) Percent {
	if goal == 0 {
		return Percent(0)
	}

	percent := float64(achieved) / float64(goal)
	return Percent(percent)
}

func (p Percent) Print() string {
	return fmt.Sprintf("%.2f", p)
}
