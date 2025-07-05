package workernotificator

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	_         = iota
	Monday    = iota
	Tuesday   = iota
	Wednesday = iota
	Thursday  = iota
	Friday    = iota
	Saturday  = iota
	Sunday    = iota
)

var WeekdaysRu = map[int]string{
	Monday:    "понедельник",
	Tuesday:   "вторник",
	Wednesday: "среда",
	Thursday:  "четверг",
	Friday:    "пятница",
	Saturday:  "суббота",
	Sunday:    "воскресенье",
}

type Notification struct {
	gorm.Model
	Name          string
	UserCategory  string // "all" || "barista" || "admin"
	WeekDays      []byte // json("[1, 4, 6]]")
	HourAndMinute string // "16:32"
	Text          string
	NowRunning    bool
}

func (n *Notification) weekDays() ([]int, error) {
	var weekDays []int
	err := json.Unmarshal(n.WeekDays, &weekDays)
	if err != nil {
		return nil, err
	}
	return weekDays, nil
}

func (n *Notification) TimeUntilNextNotification(timeZone *time.Location) (time.Duration, error) {
	now := time.Now().In(timeZone)
	currentDay := int(now.Weekday())
	currentTime := now.Hour()*60 + now.Minute()

	weekDays, err := n.weekDays()
	if err != nil {
		return 0, err
	}
	if len(weekDays) == 0 {
		return 0, errors.New("не указаны дни недели")
	}

	notificationTime, err := time.Parse("15:04", n.HourAndMinute)
	if err != nil {
		return 0, fmt.Errorf(`неправильный формат времени: "%s"`, n.HourAndMinute)
	}
	notificationMinutes := notificationTime.Hour()*60 + notificationTime.Minute()

	minDuration := time.Duration(24*7*60) * time.Minute // One week in minutes

	for _, day := range weekDays {
		if day == currentDay {
			if currentTime < notificationMinutes {
				minDuration = time.Duration(notificationMinutes-currentTime) * time.Minute
				break
			}
		} else if day > currentDay {
			daysUntilNotification := day - currentDay
			minutesUntilNotification := daysUntilNotification*24*60 + notificationMinutes - currentTime
			if time.Duration(minutesUntilNotification)*time.Minute < minDuration {
				minDuration = time.Duration(minutesUntilNotification) * time.Minute
			}
			break
		} else {
			daysUntilNotification := (7 - currentDay) + day
			minutesUntilNotification := daysUntilNotification*24*60 + notificationMinutes - currentTime
			if time.Duration(minutesUntilNotification)*time.Minute < minDuration {
				minDuration = time.Duration(minutesUntilNotification) * time.Minute
			}
		}
	}

	return minDuration, nil
}
