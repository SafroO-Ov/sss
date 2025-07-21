package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/milanakonova/dev/apifunc"
	"github.com/milanakonova/dev/db"
)

func EmployeeHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем тело запроса
	var request struct {
		EmployeeID int `json:"employee_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Получаем текущее время
	currentTime := time.Now()

	// Выводим в терминал время и ID работника
	shiftID, err := apifunc.ProcessEmployeeShift(database, request.EmployeeID, currentTime)
	if err != nil {
		log.Print("Ошибка с сменой у сотрудника")
	}
	log.Printf("Время запроса: %s, ID работника: %s, смена: %s", currentTime.Format("2006-01-02 15:04:05"), strconv.Itoa(request.EmployeeID), strconv.Itoa(shiftID))

	// Отправляем ответ клиенту
	response := map[string]string{
		"status":   "success",
		"message":  "Запрос обработан",
		"employee": strconv.Itoa(request.EmployeeID),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
