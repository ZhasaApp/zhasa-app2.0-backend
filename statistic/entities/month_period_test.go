package entities

import (
	"github.com/stretchr/testify/require"
	"testing"
	"zhasa2.0/statistic"
)

func TestConvertToTime(t *testing.T) {
	testCases := []struct {
		name         string
		monthPeriod  statistic.MonthPeriod
		expectedFrom string
		expectedTo   string
	}{
		{
			name:         "Case 1: January 2023",
			monthPeriod:  statistic.MonthPeriod{MonthNumber: 1, Year: 2023},
			expectedFrom: "01-01-2023",
			expectedTo:   "31-01-2023",
		},
		{
			name:         "Case 2: February 2023 (non-leap year)",
			monthPeriod:  statistic.MonthPeriod{MonthNumber: 2, Year: 2023},
			expectedFrom: "01-02-2023",
			expectedTo:   "28-02-2023",
		},
		{
			name:         "Case 3: February 2020 (leap year)",
			monthPeriod:  statistic.MonthPeriod{MonthNumber: 2, Year: 2020},
			expectedFrom: "01-02-2020",
			expectedTo:   "29-02-2020",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			from, to := tc.monthPeriod.ConvertToTime()
			fromStr := from.Format("02-01-2006")
			toStr := to.Format("02-01-2006")
			require.Equal(t, fromStr, tc.expectedFrom)
			require.Equal(t, toStr, tc.expectedTo)
		})
	}
}
