package dateutil

import (
	"fmt"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}

	switch {
	case repeat == "":
		return "", nil
	case repeat == "y":
		nextDate := taskDate.AddDate(1, 0, 0)
		return nextDate.Format("20060102"), nil
	case len(repeat) > 0 && repeat[0] == 'd':
		var days int
		_, err := fmt.Sscanf(repeat, "d %d", &days)
		if err != nil || days > 400 {
			return "", fmt.Errorf("invalid repeat format")
		}
		for taskDate.Before(now) {
			taskDate = taskDate.AddDate(0, 0, days)
		}
		return taskDate.Format("20060102"), nil
	default:
		return "", fmt.Errorf("unsupported repeat format")
	}
}
