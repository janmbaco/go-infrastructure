package configuration
import (
	"time"
)

type Period struct {
	weekday  time.Weekday
	hour     int
	minute   int
	isWeekly bool
}

// IsFinished returns if the period is finished
func (p *Period) IsFinished() bool {
	now := time.Now()
	return now.Hour() == p.hour && now.Minute() == p.minute && (!p.isWeekly || now.Weekday() == p.weekday)
}
