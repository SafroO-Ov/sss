package db

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func СreateNewEmployee(database *Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var newEmp NewEmployee
		if err := json.NewDecoder(r.Body).Decode(&newEmp); err != nil {
			http.Error(w, "Невалидный JSON", http.StatusBadRequest)
			return
		}

		var id int
		err := database.QueryRow(
			"INSERT INTO employees (fio, shift) VALUES ($1, $2) RETURNING employees_id",
			newEmp.FIO, newEmp.Shift,
		).Scan(&id)
		if err != nil {
			log.Printf("Ошибка вставки: %v", err)
			http.Error(w, "Ошибка при добавлении сотрудника", http.StatusInternalServerError)
			return
		}

		// Получаем созданного сотрудника для ответа
		var emp Employee
		err = database.QueryRow(
			"SELECT employees_id, fio, shift FROM employees WHERE employees_id = $1",
			id,
		).Scan(&emp.ID, &emp.FIO, &emp.Shift)
		if err != nil {
			log.Printf("Ошибка подтверждения вставки: %v", err)
			http.Error(w, "Запись не найдена после вставки", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(emp)
	}
}
