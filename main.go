package main

import (
	"log"

	"github.com/milanakonova/dev/apifunc"
	"github.com/milanakonova/dev/db"
)

func main() {
	//Начало правильного main: инициализация бд
	database, err := db.InitDB(db.NewDBConfig()) // Используем конфиг по умолчанию
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer database.Close()

	err = apifunc.StartServer("8080", database)
	if err != nil {
		log.Fatal(err)
	}

	// Правильная регистрация обработчика
	// Для добавления сотрудника: curl -X POST http://localhost:8080/employees -H "Content-Type: application/json" -d '{"fio":"Иванов Иван Иванович","shift":0}'
}
