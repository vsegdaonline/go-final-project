package api

import (
	"errors"
	"time"
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	start, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}
}
