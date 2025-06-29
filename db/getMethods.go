package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// func (r *Database) GetByIdEmployee(ctx context.Context, id int) (*Employee, error) {
func (r *Database) GetByIdEmployee(id int) (*Employee, error) {
	const query = `
		SELECT employees_id, fio, shift 
		FROM employees 
		WHERE employees_id = $1`

	var emp Employee
	err := r.QueryRow(query, id).
		Scan(&emp.ID, &emp.FIO, &emp.Shift)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("employee not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	return &emp, nil
}

// func (r *Database) GetByIdShift(ctx context.Context, id int) (*Shift, error) {
func (r *Database) GetByIdShift(id int) (*Shift, error) {
	const query = `
        SELECT shift_id, date, duration, night_time, day_time, type, on_shift, employees_id 
        FROM shifts 
        WHERE shift_id = $1`

	var (
		shift       Shift
		durationArr pq.StringArray // Используем для сканирования массива
	)

	err := r.QueryRow(query, id).
		Scan(
			&shift.ID,
			&shift.Date,
			&durationArr, // Сканируем в pq.StringArray
			&shift.NightTime,
			&shift.DayTime,
			&shift.Type,
			&shift.OnShift,
			&shift.EmployeeID,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("shift not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get shift: %w", err)
	}

	// Конвертируем pq.StringArray в []string
	shift.Duration = []string(durationArr)

	return &shift, nil
}
