package service

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCalculateRatio(t *testing.T) {
	goals := []RatioRow{
		{
			Amount:  98.0,
			Goal:    100.0,
			Gravity: 5,
		},
		{
			Amount:  29.0,
			Goal:    100.0,
			Gravity: 3,
		},
		{
			Amount:  10,
			Goal:    100,
			Gravity: 2,
		},
	}

	expected := float32(0.597)
	result := CalculateRatio(goals)

	require.Equal(t, expected, result)
}
