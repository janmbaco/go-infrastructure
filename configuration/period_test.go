package configuration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPeriod_IsFinished_WhenDailyAndTimeMatches_ThenReturnsTrue(t *testing.T) {
	// Arrange
	now := time.Now()
	period := &Period{
		hour:     now.Hour(),
		minute:   now.Minute(),
		isWeekly: false,
	}

	// Act
	finished := period.IsFinished()

	// Assert
	assert.True(t, finished)
}

func TestPeriod_IsFinished_WhenDailyAndTimeDoesNotMatch_ThenReturnsFalse(t *testing.T) {
	// Arrange
	period := &Period{
		hour:     12,
		minute:   30,
		isWeekly: false,
	}

	// Act
	finished := period.IsFinished()

	// Assert
	assert.False(t, finished)
}

func TestPeriod_IsFinished_WhenWeeklyAndTimeAndDayMatch_ThenReturnsTrue(t *testing.T) {
	// Arrange
	now := time.Now()
	period := &Period{
		weekday:  now.Weekday(),
		hour:     now.Hour(),
		minute:   now.Minute(),
		isWeekly: true,
	}

	// Act
	finished := period.IsFinished()

	// Assert
	assert.True(t, finished)
}

func TestPeriod_IsFinished_WhenWeeklyAndDayDoesNotMatch_ThenReturnsFalse(t *testing.T) {
	// Arrange
	now := time.Now()
	period := &Period{
		weekday:  now.Weekday() + 1,
		hour:     now.Hour(),
		minute:   now.Minute(),
		isWeekly: true,
	}

	// Act
	finished := period.IsFinished()

	// Assert
	assert.False(t, finished)
}
