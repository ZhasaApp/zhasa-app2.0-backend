package entities

import (
	"errors"
	"time"
)

type MonthNumber int32

func NewMonth(number int32) (MonthNumber, error) {
	if number < 1 || number > 12 {
		return 0, errors.New("month should be between 1 and 12")
	}
	return MonthNumber(number), nil
}

/*
	Period responsible for store time period in useful format and converts it into
	standard go time periods
*/
type Period interface {
	ConvertToTime() (time.Time, time.Time)
}

// MonthPeriod Period implementation for monthly period
type MonthPeriod struct {
	MonthNumber int32
	Year        int32
}

func (m *MonthPeriod) ConvertToTime() (time.Time, time.Time) {
	from := time.Date(int(m.Year), time.Month(m.MonthNumber), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)
	return from, to
}

/*
	WeekPeriod Period implementation for weekly related period
	Precondition: WeekNumber should be from 1 to 53
*/
type WeekPeriod struct {
	WeekNumber int
	Year       int
}

func (w *WeekPeriod) getMondayDate() time.Time {
	januaryFirst := time.Date(w.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	_, januaryFirstWeek := januaryFirst.ISOWeek()

	if januaryFirstWeek > 1 {
		// January 1st is in the previous year's last ISO week
		januaryFirst = time.Date(w.Year, 1, 4, 0, 0, 0, 0, time.UTC)
	}

	daysToAdd := (w.WeekNumber-1)*7 - int(januaryFirst.Weekday()-time.Monday)
	mondayDate := januaryFirst.AddDate(0, 0, daysToAdd)
	return mondayDate
}

func (w *WeekPeriod) ConvertToTime() (time.Time, time.Time) {
	from := w.getMondayDate()
	to := from.AddDate(0, 0, 6)
	return from, to
}

type DayPeriod struct {
	Day time.Time
}

func (d DayPeriod) ConvertToTime() (time.Time, time.Time) {
	from := time.Date(d.Day.Year(), d.Day.Month(), d.Day.Day(), 0, 0, 0, 0, d.Day.Location())
	to := from.Add(time.Hour * 24).Add(-time.Nanosecond)
	return from, to
}
