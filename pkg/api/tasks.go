package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"go1f/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		err := errors.New("не указан идентификатор")
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		err := errors.New("задача не найдена")
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	task := db.Task{}
	var buf bytes.Buffer
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		err = errors.New("не указан заголовок задачи")
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err = checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		err = errors.New("задача не найдена")
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{})
}
