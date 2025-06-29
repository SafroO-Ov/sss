package apifunc

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/milanakonova/dev/db"
	"github.com/milanakonova/dev/handlers"
)

// StartServer запускает HTTP сервер на указанном порту
func StartServer(port string, database *db.Database) error {
	// Регистрируем обработчики с передачей db через замыкание
	http.HandleFunc("/employee", func(w http.ResponseWriter, r *http.Request) {
		employeeHandler(w, r, database)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
	})
	http.HandleFunc("/employee-stats", func(w http.ResponseWriter, r *http.Request) {
		handlers.EmployeeStatsHandler(w, r, database)
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}
	})
	http.HandleFunc("/newemployee", СreateNewEmployee(database))
	http.HandleFunc("/allemployees", handlers.GetAllEmployeesHandler(database))

	addr := ":" + port
	log.Printf("Сервер запущен на http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}

// Обработчик для пути /employee
func employeeHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
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
	shiftID, err := ProcessEmployeeShift(database, request.EmployeeID, currentTime)
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
