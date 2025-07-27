package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/milanakonova/dev/db"
)

// EmployeeStatsRequest структура для входящего запроса

// EmployeeStatsResponse структура для ответа
type EmployeeStatsResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Stats      interface{} `json:"stats"`
	EmployeeID int         `json:"employee_id"`
	Mode       string      `json:"mode"`
}

func EmployeeStatsHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Для предварительных OPTIONS запросов
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Декодирование тела запроса
	var request db.EmployeeStatsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Валидация режима
	if request.Mode != "current" && request.Mode != "week" && request.Mode != "month" {
		http.Error(w, "Некорректный режим. Допустимые значения: current, week, month", http.StatusBadRequest)
		return
	}
	if request.Mode == "current" {
		stats, err := db.GetCurrentShiftStats(database, request.EmployeeID)
		if err != nil {
			log.Printf("Ошибка при получении статистики: %v", err)
			http.Error(w, "Ошибка при обработке запроса", http.StatusInternalServerError)
			return
		}
		fmt.Println(stats)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	} else {
		stats, err := db.GetShiftStatsMore(database, request.EmployeeID, request.Mode)
		if err != nil {
			log.Printf("Ошибка при получении статистики: %v", err)
			http.Error(w, "Ошибка при обработке запроса", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
