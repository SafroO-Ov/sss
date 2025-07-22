package db

// GetAllEmployees возвращает всех сотрудников из базы данных
func GetAllEmployees(database *Database) ([]Employee, error) {
	query := `SELECT employees_id, fio, shift FROM employees`
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.ID, &emp.FIO, &emp.Shift); err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}
