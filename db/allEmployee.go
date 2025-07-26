package db

import "log"

// GetAllEmployees возвращает всех сотрудников из базы данных
func GetAllEmployees(database *Database) ([]EmployeeWithShiftStatus, error) {
	query := `
        SELECT e.employees_id, e.fio, e.shift, s.on_shift 
        FROM employees e
        LEFT JOIN shifts s ON e.shift = s.shift_id
    `
	rows, err := database.Query(query)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v\n", err) // Логируем ошибку
		return nil, err
	}
	defer rows.Close()

	var employees []EmployeeWithShiftStatus
	for rows.Next() {
		var emp EmployeeWithShiftStatus
		if err := rows.Scan(
			&emp.ID,
			&emp.FIO,
			&emp.Shift,
			&emp.OnShift,
		); err != nil {
			log.Printf("Ошибка при сканировании строки: %v\n", err) // Логируем ошибку
			return nil, err
		}
		employees = append(employees, emp)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Ошибка после обработки строк: %v\n", err) // Логируем ошибку
		return nil, err
	}

	// Выводим полученные данные в терминал
	log.Println("Получены сотрудники со статусом смены:")
	for i, emp := range employees {
		log.Printf(
			"[%d] ID: %d, ФИО: %s, ID смены: %v, На смене: %v\n",
			i+1,
			emp.ID,
			emp.FIO,
			emp.Shift,
			emp.OnShift,
		)
	}

	return employees, nil
}
