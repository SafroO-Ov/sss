package db

import (
	"database/sql"

	"github.com/lib/pq"
)

type Employee struct {
	ID    int    `json:"employees_id" db:"employees_id"`
	FIO   string `json:"fio" db:"fio"`
	Shift int    `json:"shift" db:"shift"`
}

type NewEmployee struct {
	FIO   string `json:"fio" validate:"required"`
	Shift int    `json:"shift" validate:"required"`
}
type Shift struct {
	ID         int            `json:"shift_id" db:"shift_id"`
	Date       string         `json:"date" db:"date"`
	Duration   pq.StringArray `json:"duration" db:"duration"`
	NightTime  int            `json:"night_time" db:"night_time"`
	DayTime    int            `json:"day_time" db:"day_time"`
	Type       string         `json:"type" db:"type"`
	OnShift    bool           `json:"on_shift" db:"on_shift"`
	EmployeeID int            `json:"employees_id" db:"employees_id"`
}

type NewShift struct {
	Date       string   `json:"date" validate:"required"`
	Duration   []string `json:"duration" validate:"required"`
	Type       string   `json:"type" validate:"required"`
	EmployeeID int      `json:"employees_id" validate:"required"`
}

// Database структура для работы с БД
type Database struct {
	*sql.DB
}

// DBConfig конфигурация подключения к БД
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}
type EmployeeStatsRequest struct {
	EmployeeID int    `json:"employee_id"`
	Mode       string `json:"mode"` // "current", "week" или "month"
}
type ShiftStats struct {
	DayHours   int `json:"day_hours"`
	NightHours int `json:"night_hours"`
	TotalHours int `json:"total_hours"`
}
type ShiftStatsMore struct {
	ShiftCount  int `json:"shift_count"`
	NightShifts int `json:"night_shifts"`
	TotalHours  int `json:"total_hours"`
	NightHours  int `json:"night_hours"`
	OverTime    int `json:"overtime_hours"`
}
type EmployeeWithShiftStatus struct {
	Employee      // Встраиваем оригинальную структуру
	OnShift  bool // или bool, если NULL не ожидается
}
