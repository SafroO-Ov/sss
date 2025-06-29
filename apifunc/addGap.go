package apifunc

import (
	"fmt"
	"log"
	"time"

	"github.com/milanakonova/dev/db"
)

// ProcessEmployeeShift обрабатывает смену сотрудника
func ProcessEmployeeShift(database *db.Database, employeeID int, requestTime time.Time) (shiftID int, err error) {
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
			err = updateShiftDuration(database, sh)
			if err != nil {
				return 0, err
			}

		} else {
			// Четное количество - проверяем временной промежуток
			lastTime, err := time.Parse("15:04:05", sh.Duration[len(sh.Duration)-1])
			if err != nil {
				return 0, fmt.Errorf("ошибка парсинга времени: %w", err)
			}

			currentTime, _ := time.Parse("15:04:05", requestTime.Format("15:04:05"))
			hoursPassed := currentTime.Sub(lastTime).Hours()
			// добавить чтоб хранилось и считалось в секундах
			if hoursPassed >= 7 {
				// Создаем новую смену
				newsh, err := CreateNewShift(database, employeeID, requestTime)
				if err != nil {
					fmt.Println(err)
					return 0, err
				}
				return newsh.ID, err
			} else {
				// Добавляем текущее время в duration
				sh.Duration = append(sh.Duration, requestTime.Format("15:04:05"))

				err := updateShiftDuration(database, sh)
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
