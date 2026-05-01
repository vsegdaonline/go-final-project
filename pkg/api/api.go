package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"go1f/pkg/db"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

const dateFormat = "20060102"

func Init(r *chi.Mux) {
	r.Post("/api/signin", signinHandler)
	r.Get("/api/nextdate", auth(nextDayHandler))
	r.Post("/api/task", auth(addTaskHandler))
	r.Get("/api/tasks", auth(tasksHandler))
	r.Get("/api/task", auth(getTaskHandler))
	r.Put("/api/task", auth(putTaskHandler))
	r.Post("/api/task/done", auth(doneTaskHandler))
	r.Delete("/api/task", auth(deleteTaskHandler))
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	passNeed := os.Getenv("TODO_PASSWORD")
	if passNeed == "" {
		writeJson(w, map[string]string{})
		return
	}
	type frontPass struct {
		Password string `json:"password"`
	}
	passGet := frontPass{}
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &passGet); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	if passNeed != passGet.Password {
		err = errors.New("неверный пароль")
		w.WriteHeader(http.StatusUnauthorized)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	secret := []byte(passNeed)
	hash := sha256.Sum256([]byte(passNeed))
	hashString := hex.EncodeToString(hash[:])
	claims := jwt.MapClaims{
		"hash": hashString,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{"token": signedToken})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtS string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtS = cookie.Value
			} else {
				http.Error(w, "cookie is empty", http.StatusUnauthorized)
				return
			}
			var valid bool
			token, err := jwt.Parse(jwtS, func(t *jwt.Token) (interface{}, error) {
				return []byte(pass), nil
			})
			valid = token.Valid
			if err != nil || !valid {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "cookie is not correct", http.StatusUnauthorized)
				return
			}
			hashFromToken, ok := claims["hash"].(string)
			if !ok {
				http.Error(w, "cookie is not correct", http.StatusUnauthorized)
				return
			}
			hash := sha256.Sum256([]byte(pass))
			hashString := hex.EncodeToString(hash[:])
			if hashFromToken != hashString {
				http.Error(w, "cookie is not correct", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

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

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	//сразу устанавливаю заголовок
	//так как и ошибки и успешные ответы будут в json формате
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

	id, err := db.AddTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{"id": strconv.Itoa(int(id))})
}

func checkDate(task *db.Task) error {
	nowStr := time.Now().Format(dateFormat)
	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		return err
	}
	if task.Date == "" {
		task.Date = nowStr
	}
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = nowStr
		} else {
			task.Date = next
		}
	}
	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	_ = json.NewEncoder(w).Encode(data)
}
