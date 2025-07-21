package db

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

// CreateNewShift создаёт новую смену для сотрудника CreateNewShift(db, employeeID, requestTime)
func CreateNewShift(database *Database, employeeID int, requestTime time.Time) (*Shift, error) {
	// Проверяем валидность подключения к БД
	if database == nil {
		return nil, fmt.Errorf("невалидное подключение к БД")
	}

	// Создаем объект новой смены
	newShift := NewShift{
		Date:       requestTime.Format("2006-01-02"),
		Duration:   []string{requestTime.Format("15:04:05")},
		Type:       "дневная",
		EmployeeID: employeeID,
	}

	// Выполняем запрос к базе данных
	query := `INSERT INTO shifts 
		(date, duration, night_time, day_time, type, on_shift, employees_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING shift_id, date, duration, night_time, day_time, type, on_shift, employees_id`

	var (
		createdShift Shift
		durationArr  pq.StringArray // Используем для сканирования массива
	)
	err := database.QueryRow(
		query,
		newShift.Date,
		pq.Array(newShift.Duration),
		0, // night_time
		0, // day_time
		newShift.Type,
		true, // on_shift
		newShift.EmployeeID,
	).Scan(
		&createdShift.ID,
		&createdShift.Date,
		&durationArr, // Сканируем в pq.StringArray
		&createdShift.NightTime,
		&createdShift.DayTime,
		&createdShift.Type,
		&createdShift.OnShift,
		&createdShift.EmployeeID,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка при создании смены: %w", err)
	}

	// Конвертируем pq.StringArray в []string
	createdShift.Duration = durationArr

	// Обновляем смену у сотрудника
	updateQuery := "UPDATE employees SET shift = $1 WHERE employees_id = $2"
	_, err = database.Exec(updateQuery, createdShift.ID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении сотрудника: %w", err)
	}

	return &createdShift, nil
}
