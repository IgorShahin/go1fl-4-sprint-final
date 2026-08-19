package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	slicesData := strings.Split(data, ",")
	if len(slicesData) != 2 {
		return 0, 0, fmt.Errorf("expected 2 items, got %d", len(slicesData))
	}

	steps, err := strconv.Atoi(slicesData[0])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("invalid number of steps: %d", steps)
	}

	timeData, err := time.ParseDuration(slicesData[1])
	if err != nil {
		return 0, 0, err
	}

	if timeData <= 0 {
		return 0, 0, fmt.Errorf("invalid time: %s", timeData)
	}

	return steps, timeData, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, timeData, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	distance := (float64(steps) * stepLength) / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, timeData)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps,
		distance,
		calories,
	)

	return result
}
