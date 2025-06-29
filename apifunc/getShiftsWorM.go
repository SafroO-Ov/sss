package apifunc

import (
	"time"

	"github.com/lib/pq"
	"github.com/milanakonova/dev/db"
)

// GetEmployeeShiftsForWeek возвращает все смены сотрудника за прошедшую неделю
func GetEmployeeShiftsForWeek(database *db.Database, employeeID int) ([]db.Shift, error) {
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)

	rows, err := database.Query(`
		SELECT shift_id, date, duration, night_time, day_time, type, on_shift, employees_id 
		FROM shifts 
		WHERE employees_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC`,
		employeeID, weekAgo.Format("2006-01-02"), now.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []db.Shift
	for rows.Next() {
		var shift db.Shift
		err := rows.Scan(
			&shift.ID,
			&shift.Date,
			pq.Array(&shift.Duration),
			&shift.NightTime,
			&shift.DayTime,
			&shift.Type,
			&shift.OnShift,
			&shift.EmployeeID,
		)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, shift)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shifts, nil
}

// GetEmployeeShiftsForMonth возвращает все смены сотрудника за прошедший месяц
func GetEmployeeShiftsForMonth(database *db.Database, employeeID int) ([]db.Shift, error) {
	now := time.Now()
	monthAgo := now.AddDate(0, -1, 0)

	rows, err := database.Query(`
		SELECT shift_id, date, duration, night_time, day_time, type, on_shift, employees_id 
		FROM shifts 
		WHERE employees_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC`,
		employeeID, monthAgo.Format("2006-01-02"), now.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []db.Shift
	for rows.Next() {
		var shift db.Shift
		err := rows.Scan(
			&shift.ID,
			&shift.Date,
			pq.Array(&shift.Duration),
			&shift.NightTime,
			&shift.DayTime,
			&shift.Type,
			&shift.OnShift,
			&shift.EmployeeID,
		)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, shift)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shifts, nil
}
