package db

import (
	"database/sql"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int, search string) ([]*Task, error) {
	tasks := []*Task{}
	if search != "" {
		tasks, err := searchTasks(limit, search)
		if err != nil {
			return tasks, err
		}
		return tasks, nil
	}
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`
	tasksRows, err := db.Query(query, sql.Named("limit", limit))
	if err != nil {
		return tasks, err
	}
	defer func() {
		_ = tasksRows.Close()
	}()
	for tasksRows.Next() {
		var task Task
		err = tasksRows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &task)
	}
	if err = tasksRows.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}

func searchTasks(limit int, search string) ([]*Task, error) {
	tasks := []*Task{}
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`
	queryByDate := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit`

	date, ok := isDate(search)
	if ok {
		tasksRows, err := db.Query(queryByDate, sql.Named("date", date), sql.Named("limit", limit))
		if err != nil {
			return tasks, err
		}
		defer func() { _ = tasksRows.Close() }()

		for tasksRows.Next() {
			var task Task
			err = tasksRows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return tasks, err
			}
			tasks = append(tasks, &task)
		}
		if err = tasksRows.Err(); err != nil {
			return tasks, err
		}
		return tasks, nil
	}
	tasksRows, err := db.Query(query, sql.Named("search", "%"+search+"%"), sql.Named("limit", limit))
	if err != nil {
		return tasks, err
	}
	defer func() { _ = tasksRows.Close() }()
	for tasksRows.Next() {
		var task Task
		err = tasksRows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &task)
	}
	if err = tasksRows.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}

func isDate(search string) (string, bool) {
	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		return "", false
	}
	dateStr := date.Format("20060102")
	return dateStr, true
}
