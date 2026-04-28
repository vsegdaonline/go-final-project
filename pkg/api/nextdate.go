package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}
	repeatSlice := strings.Split(repeat, " ")
	switch repeatSlice[0] {
	case "y":
		if len(repeatSlice) != 1 {
			return "", errors.New("repeat have not correct format")
		}
		for {
			date = date.AddDate(1, 0, 0)

			if afterNow(date, now) {
				break
			}
		}
	case "m":
		if len(repeatSlice) > 3 || len(repeatSlice) < 2 {
			return "", errors.New("repeat have not correct format")
		}
		var day [33]bool
		var month [13]bool
		fPartRepeat := strings.Split(repeatSlice[1], ",")
		for _, v := range fPartRepeat {
			dayNumber, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			if dayNumber > 31 {
				return "", errors.New("недопустимый день месяца")
			}
			if dayNumber <= 0 {
				switch dayNumber {
				case -1:
					day[0] = true
				case -2:
					day[32] = true
				default:
					return "", errors.New("недопустимый день месяца")
				}
			} else {
				day[dayNumber] = true
			}
		}
		if len(repeatSlice) == 3 {
			sPartRepeat := strings.Split(repeatSlice[2], ",")
			for _, v := range sPartRepeat {
				monthNumber, err := strconv.Atoi(v)
				if err != nil {
					return "", err
				}
				if monthNumber < 1 || monthNumber > 12 {
					return "", errors.New("недопустимый месяц")
				}
				month[monthNumber] = true
			}
		} else {
			for i := range month {
				if i == 0 {
					continue
				}
				month[i] = true
			}
		}
		ok := true
		count := 0
		for ok {
			date = date.AddDate(0, 0, 1)
			if month[date.Month()] {
				switch {
				case day[date.Day()]:
					if afterNow(date, now) {
						ok = false
						break
					}
				case day[0]:
					if date.AddDate(0, 0, 1).Month() != date.Month() {
						if afterNow(date, now) {
							ok = false
							break
						}
					}
				case day[32]:
					if date.AddDate(0, 0, 2).Month() != date.Month() && date.AddDate(0, 0, 1).Month() == date.Month() {
						if afterNow(date, now) {
							ok = false
							break
						}
					}
				}
			}
			if count > 1500 {
				return "", errors.New("недопустимый день месяца")
			}
			count++
		}
	case "w":
		var week [7]bool
		if len(repeatSlice) != 2 {
			return "", errors.New("repeat have not correct format")
		}
		weekSlice := strings.Split(repeatSlice[1], ",")
		for _, v := range weekSlice {
			day, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			if day < 1 || day > 7 {
				return "", errors.New("недопустимое значение дня недели")
			}
			if day == 7 {
				week[0] = true
			} else {
				week[day] = true
			}
		}
		for {
			date = date.AddDate(0, 0, 1)
			if week[date.Weekday()] {
				if afterNow(date, now) {
					break
				}
			}
		}
	case "d":
		if len(repeatSlice) != 2 {
			return "", errors.New("repeat have not correct format")
		}
		day, err := strconv.Atoi(repeatSlice[1])
		if err != nil {
			return "", err
		}
		if day > 400 || day < 1 {
			return "", errors.New("day is out of range 1-400")
		}
		for {
			date = date.AddDate(0, 0, day)
			if afterNow(date, now) {
				break
			}
		}
	default:
		return "", errors.New("symbol is not correct")
	}

	return date.Format("20060102"), nil
}
