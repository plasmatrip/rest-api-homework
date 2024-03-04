package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// isEmpty возвращает true, если мапа с задачами пустая, иначе возвращает false
// если мапа пустая возвращет клиенту в ответ ошибку
func isEmpty(w http.ResponseWriter) bool {
	if len(tasks) == 0 {
		http.Error(w, "Список задач пустой", http.StatusBadRequest)
		return true
	}
	return false
}

// getTasks возвращает все задачи из мапы
// обработчик для маршрута `/tasks` с методом GET
func getTasks(w http.ResponseWriter, r *http.Request) {
	//проверяем пустая мапа или нет
	if isEmpty(w) {
		return
	}

	//сериализуем данные из мапы
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//записываем в заголовок тип контента
	w.Header().Set("Content-Type", "application/json")
	//статус ответа
	w.WriteHeader(http.StatusOK)
	//записывем сериализованные данные в тело ответа
	w.Write(resp)
}

// getTask возвращает задачу из мапы по id
// обработчик для маршрута `/tasks/{id}` с методом GET
func getTask(w http.ResponseWriter, r *http.Request) {
	//проверяем пустая мапа или нет
	if isEmpty(w) {
		return
	}

	//читаем данные из мапы по id
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	//если данные не найдены, передаем клиенту ошибку
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	//сериализуем найденные данные
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//записываем в заголовок тип контента
	w.Header().Set("Content-Type", "application/json")
	//статус ответа
	w.WriteHeader(http.StatusOK)
	//записывем сериализованные данные в тело ответа
	w.Write(resp)
}

// postTask добавляет в мапу новую задачу
// обработчик для маршрута `/tasks` с методом POST
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	//читаем из тела запроса данные в буфер
	_, err := buf.ReadFrom(r.Body)
	//если возникла ошибка, передаем ее клиенту
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//сериализуем прочитанные данные в переменную task, если ошибка - передаем ее клиенту
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//если задача с переданным id есть, возвращаем ошибку
	_, ok := tasks[task.ID]
	if ok {
		errStr := fmt.Sprintf("Задача c id=%s уже есть", task.ID)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	//записываем задачу в мапу
	tasks[task.ID] = task

	//записываем в заголовок тип контента
	w.Header().Set("Content-Type", "application/json")
	//статус ответа
	w.WriteHeader(http.StatusCreated)
}

// delTask удаляет задачу из мапы по id
// обработчик для маршрута `/tasks` с методом DELETE
func delTask(w http.ResponseWriter, r *http.Request) {
	//проверяем пустая мапа или нет
	if isEmpty(w) {
		return
	}

	//читаем параметр id из URL
	id := chi.URLParam(r, "id")
	//ищем задачу в мапе по id
	_, ok := tasks[id]
	//если задача не найдена, возвращаем клиенту ошибку
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	//удаляем задачу из мапы
	delete(tasks, id)

	//записываем в заголовок тип контента
	w.Header().Set("Content-Type", "application/json")
	//статус ответа
	w.WriteHeader(http.StatusOK)
}

func main() {
	//создаем новый ройтер
	r := chi.NewRouter()

	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTasks`
	r.Get("/tasks", getTasks)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTask`
	r.Get("/tasks/{id}", getTask)
	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTask`
	r.Post("/tasks", postTask)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `delTask`
	r.Delete("/tasks/{id}", delTask)

	//запускаем сервер
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
