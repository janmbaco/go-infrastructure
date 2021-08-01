package configuration

import (
	"time"
)

// Period defines an object responsible for knowing if the period is finished
type Period interface {
	IsFinished() bool
}

type period struct {
	weekday  time.Weekday
	hour     int
	minute   int
	isWeekly bool
}

func (period *period) IsFinished() bool {
	now := time.Now()
	return now.Hour() == period.hour && now.Minute() == period.minute && (!period.isWeekly || now.Weekday() == period.weekday)
}

// NewWeeklyPeriod returns a weekly period
func NewWeeklyPeriod(weekDays time.Weekday, hour int, minute int) Period {
	return &period{weekday: weekDays, hour: hour, minute: minute, isWeekly: true}
}

// NewDailyPeriod returns a daily period
func NewDailyPeriod(hour int, minute int) Period {
	return &period{hour: hour, minute: minute}
}
