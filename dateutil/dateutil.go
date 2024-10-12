package dateutil

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату на основе указанного правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("не указано правило повторения")
	}

	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	// Проверяем правило повторения
	switch {
	case repeat == "y":
		nextDate := startDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0) // продолжаем увеличивать до нахождения подходящей даты
		}
		return nextDate.Format("20060102"), nil

	case strings.HasPrefix(repeat, "d "):
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("неверное значение дней в правиле d")
		}

		nextDate := startDate
		// Продолжаем добавлять дни, пока дата не станет больше now
		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if nextDate.After(now) {
				return nextDate.Format("20060102"), nil
			}
		}

	case strings.HasPrefix(repeat, "m "):
		parts := strings.Split(repeat, " ")
		daysStr := parts[1]
		var days []int

		dayParts := strings.Split(daysStr, ",")
		for _, day := range dayParts {
			dayInt, err := strconv.Atoi(day)
			if err != nil {
				return "", errors.New("неверное значение дня месяца")
			}
			days = append(days, dayInt)
		}

		nextDate := startDate
		for {
			if nextDate.After(now) {
				break
			}
			nextDate = nextDate.AddDate(0, 0, 1) // Увеличиваем день
			for _, day := range days {
				if day == -1 {
					nextDate = time.Date(nextDate.Year(), nextDate.Month(), 1, 0, 0, 0, 0, nextDate.Location())
					nextDate = nextDate.AddDate(0, 1, -1) // последний день месяца
				} else if day == -2 {
					nextDate = time.Date(nextDate.Year(), nextDate.Month(), 1, 0, 0, 0, 0, nextDate.Location())
					nextDate = nextDate.AddDate(0, 1, -2) // предпоследний день месяца
				} else if nextDate.Day() == day {
					return nextDate.Format("20060102"), nil
				}
			}
		}
		return nextDate.Format("20060102"), nil

	case strings.HasPrefix(repeat, "w "):
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила w")
		}
		dayParts := strings.Split(parts[1], ",")
		var weekDays []int

		for _, day := range dayParts {
			dayInt, err := strconv.Atoi(day)
			if err != nil || dayInt < 1 || dayInt > 7 {
				return "", errors.New("неверное значение дня недели")
			}
			weekDays = append(weekDays, dayInt)
		}

		nextDate := startDate
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			if nextDate.After(now) {
				for _, day := range weekDays {
					if int(nextDate.Weekday()) == (day % 7) { // Sunday is 0 in Go
						return nextDate.Format("20060102"), nil
					}
				}
			}
		}

	default:
		return "", errors.New("неподдерживаемый формат")
	}

	//return "", errors.New("не удалось найти следующую дату")
}
