package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func GetCurrentShiftStats(database *Database, employeeID int) (*ShiftStats, error) {
	// 1. Получаем ID смены сотрудника
	var shiftID int
	err := database.QueryRow(`
		SELECT shift 
		FROM employees 
		WHERE employees_id = $1`, employeeID).Scan(&shiftID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("сотрудник не найден")
		}
		return nil, fmt.Errorf("ошибка получения смены: %w", err)
	}

	// 2. Получаем данные смены
	var dayTime, nightTime int
	err = database.QueryRow(`
		SELECT day_time, night_time 
		FROM shifts 
		WHERE shift_id = $1`, shiftID).Scan(&dayTime, &nightTime)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("смена не найдена")
		}
		return nil, fmt.Errorf("ошибка получения данных смены: %w", err)
	}

	// 3. Формируем результат
	return &ShiftStats{
		DayHours:   dayTime,
		NightHours: nightTime,
		TotalHours: dayTime + nightTime,
	}, nil
}
func GetShiftStatsMore(database *Database, employeeID int, periodType string) (*ShiftStatsMore, error) {
	// Проверяем валидность периода
	var timeRange string
	switch periodType {
	case "week":
		timeRange = "7 days"
	case "month":
		timeRange = "30 days"
	default:
		return nil, errors.New("неподдерживаемый тип периода (допустимо: week, month)")
	}

	// Проверяем существование сотрудника
	var exists bool
	err := database.QueryRow("SELECT EXISTS(SELECT 1 FROM employees WHERE employees_id = $1)", employeeID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки сотрудника: %w", err)
	}
	if !exists {
		return nil, errors.New("сотрудник не найден")
	}

	// Выполняем запрос для получения статистики
	query := `
		SELECT 
			COUNT(*) as shift_count,
			SUM(CASE WHEN s.type = 'night' THEN 1 ELSE 0 END) as night_shifts,
			COALESCE(SUM(s.day_time + s.night_time), 0) as total_hours,
			COALESCE(SUM(s.night_time), 0) as night_hours,
			COALESCE(SUM(s.day_time), 0) as day_hours,
			COALESCE(SUM(CASE 
				WHEN (s.day_time + s.night_time) > 8 THEN (s.day_time + s.night_time) - 8 
				ELSE 0 
			END), 0) as overtime_hours
		FROM shifts s
		JOIN employees e ON s.employee_id = e.employees_id
		WHERE e.employees_id = $1
		AND s.date BETWEEN NOW() - $2::INTERVAL AND NOW()
	`

	var stats ShiftStatsMore
	err = database.QueryRow(query, employeeID, timeRange).Scan(
		&stats.ShiftCount,
		&stats.NightShifts,
		&stats.TotalHours,
		&stats.NightHours,
		&stats.OverTime,
	)
	if timeRange == "7 days" {
		stats.OverTime -= 40
	} else {
		stats.OverTime -= 160
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения статистики: %w", err)
	}

	return &stats, nil
}
