package db

import (
	"fmt"
	"log"
	"time"
)

// ProcessEmployeeShift обрабатывает смену сотрудника
func ProcessEmployeeShift(database *Database, employeeID int, requestTime time.Time) (shiftID int, err error) {
	// 1. Получаем данные сотрудника из БД
	empl, err := database.GetByIdEmployee(employeeID)
	if err != nil {
		log.Println("Ошибка при запросе1")
	}

	// 2. Если смена назначена - создаем новую
	// если чел уходит со смены то должна вызваться функция которая посчитает часов за смену
	if empl.Shift != 0 {
		sh, err := database.GetByIdShift(empl.Shift)
		if err != nil {
			log.Println("Ошибка при запросе2 ", err)
		}
		if len(sh.Duration)%2 == 1 {
			// Нечетное количество - добавляем текущее время
			sh.Duration = append(sh.Duration, requestTime.Format("15:04:05"))

			// Обновляем смену в БД
			err = UpdateShiftDuration(database, sh)
			if err != nil {
				return 0, err
			}

		} else {
			// Четное количество - проверяем временной промежуток

			// Парсим последнее время из Duration с учетом даты из sh.Date
			lastDateTimeStr := sh.Date + " " + sh.Duration[len(sh.Duration)-1]
			lastDateTime, err := time.Parse("2006-01-02 15:04:05", lastDateTimeStr)
			if err != nil {
				return 0, fmt.Errorf("ошибка парсинга времени: %w", err)
			}

			// Получаем разницу между текущим временем и последним временем в смене
			timeDiff := requestTime.Sub(lastDateTime)
			hoursPassed := timeDiff.Hours()

			// Если прошло 7 или более часов
			if hoursPassed >= 7 {
				// Создаем новую смену
				newsh, err := CreateNewShift(database, employeeID, requestTime)
				if err != nil {
					fmt.Println(err)
					return 0, err
				}
				return newsh.ID, err
			} else {
				// Добавляем текущее время в duration (только время без даты)
				sh.Duration = append(sh.Duration, requestTime.Format("15:04:05"))

				err := UpdateShiftDuration(database, sh)
				if err != nil {
					fmt.Println(err)
					return 0, err
				}
			}
		}
		return sh.ID, err

	} else {
		newshift, err := CreateNewShift(database, employeeID, requestTime)
		if err != nil {
			log.Println("Ошибка при запросе3")
			return 0, err
		}
		return newshift.ID, err
	}
}
