package db

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

func UpdateShiftDuration(database *Database, shift *Shift) error {
	// Проверяем количество элементов в duration
	if len(shift.Duration)%2 == 0 {
		// Четное количество - завершаем смену
		shift.OnShift = false
		startTime, err1 := time.Parse("15:04:05", shift.Duration[len(shift.Duration)-2])
		endTime, err2 := time.Parse("15:04:05", shift.Duration[len(shift.Duration)-1])

		if err1 == nil && err2 == nil {
			// Рассчитываем ночные и дневные часы
			nightHours, dayHours := CalculateShiftHours(startTime, endTime)
			fmt.Println(shift.NightTime, shift.DayTime)
			// Обновляем суммарное время
			shift.NightTime += nightHours
			shift.DayTime += dayHours
			if shift.NightTime > 0 && shift.DayTime == 0 {
				shift.Type = "ночная"
			}
			if shift.NightTime > 0 && shift.DayTime > 0 {
				shift.Type = "гибридная"
			}

			// Обновляем дату смены на текущую дату
			shift.Date = time.Now().Format("2006-01-02")
		}

	} else {
		// Нечетное количество - смена активна
		shift.OnShift = true
	}

	// Обновляем запись в базе данных
	_, err := database.Exec(
		`UPDATE shifts 
         SET duration = $1, on_shift = $2, night_time = $3, 
             day_time = $4, type = $5, date = $6
         WHERE shift_id = $7`,
		pq.Array(shift.Duration),
		shift.OnShift,
		shift.NightTime,
		shift.DayTime,
		shift.Type,
		shift.Date, // Добавлено обновление даты
		shift.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления смены: %w", err)
	}

	return nil
}

// CalculateShiftHours рассчитывает ночные и дневные часы между двумя временными точками
func CalculateShiftHours(startTime, endTime time.Time) (nightHours, dayHours int) {
	const (
		nightStart = 22 // 22:00
		nightEnd   = 7  // 07:00
	)

	if endTime.Before(startTime) {
		endTime = endTime.Add(24 * time.Hour)
	}

	current := startTime

	for current.Before(endTime) {
		next := current.Add(time.Minute)
		if next.After(endTime) {
			next = endTime
		}

		hour := current.Hour()

		if hour >= nightStart || hour < nightEnd {
			nightHours++
		} else {
			dayHours++
		}

		current = next
	}

	// Минуты → часы округляем вниз
	return nightHours / 60, dayHours / 60
}
