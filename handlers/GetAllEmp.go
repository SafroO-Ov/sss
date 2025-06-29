package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/milanakonova/dev/db"
)

func GetAllEmployeesHandler(database *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}
		employees, err := db.GetAllEmployees(database)
		if err != nil {
			log.Printf("ошибка получения сотрудников: %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(employees); err != nil {
			log.Printf("ошибка кодирования JSON: %v", err)
			http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
		}
	}
}
