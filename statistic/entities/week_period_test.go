package entities

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"zhasa2.0/statistic"
)

func TestGetMondayDate(t *testing.T) {
	testCases := []struct {
		name     string
		w        statistic.WeekPeriod
		expected time.Time
	}{
		{
			name: "Case 1: Week 20 in 2023",
			w: statistic.WeekPeriod{
				Year:       2023,
				WeekNumber: 20,
			},
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Case 2: Week 1 in 2022",
			w: statistic.WeekPeriod{
				Year:       2022,
				WeekNumber: 1,
			},
			expected: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Case 3: Week 53 in 2020",
			w: statistic.WeekPeriod{
				Year:       2020,
				WeekNumber: 53,
			},
			expected: time.Date(2020, 12, 28, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			from := tc.w.GetMondayDate()
			require.Equal(t, from, tc.expected)
		})
	}
}
