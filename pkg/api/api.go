package api

import (
	"net/http"
	"time"
)

const dateFormat = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if now == "" {
		now = time.Now().Format(dateFormat)
	}
	nowTime, err := time.Parse(dateFormat, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(result))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
}
