package apifunc

import (
	"log"
	"net/http"

	"github.com/milanakonova/dev/db"
	"github.com/milanakonova/dev/handlers"
)

// StartServer запускает HTTP сервер на указанном порту
func StartServer(port string, database *db.Database) error {
	// Регистрируем обработчики с передачей db через замыкание
	http.HandleFunc("/employee", func(w http.ResponseWriter, r *http.Request) {
		handlers.EmployeeHandler(w, r, database)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	})
	http.HandleFunc("/employee-stats", func(w http.ResponseWriter, r *http.Request) {
		handlers.EmployeeStatsHandler(w, r, database)
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}
	})
	http.HandleFunc("/newemployee", db.СreateNewEmployee(database))
	http.HandleFunc("/allemployees", handlers.GetAllEmployeesHandler(database))

	addr := ":" + port
	log.Printf("Сервер запущен на http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}
