package base

import (
	"fmt"
	"math"
)

type Percent float64

func NewPercent(achieved int64, goal int64) Percent {
	if goal == 0 {
		return Percent(0)
	}

	percent := float64(achieved) / float64(goal)
	roundedPercent := math.Round(percent*100) / 100
	return Percent(roundedPercent * 100)
}

func (p Percent) GetRounded() Percent {
	roundedPercent := math.Round(float64(p)*100) / 100
	return Percent(roundedPercent * 100)
}

func (p Percent) Print() string {
	return fmt.Sprintf("%.2f", p)
}
