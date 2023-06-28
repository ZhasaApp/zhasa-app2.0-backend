package service

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCalculateRatio(t *testing.T) {
	goals := []RatioRow{
		{
			amount:  80.0,
			goal:    100.0,
			gravity: 5,
		},
		{
			amount:  70.0,
			goal:    100.0,
			gravity: 3,
		},
		{
			amount:  40,
			goal:    100,
			gravity: 2,
		},
	}

	expected := float32(0.69)
	result := CalculateRatio(goals)

	require.Equal(t, expected, result)
}