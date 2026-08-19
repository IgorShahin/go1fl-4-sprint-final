package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	slicesData := strings.Split(data, ",")
	if len(slicesData) != 3 {
		return 0, "", 0, fmt.Errorf("expected 3 items, got %d", len(slicesData))
	}

	steps, err := strconv.Atoi(slicesData[0])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("invalid number of steps: %d", steps)
	}

	activity := slicesData[1]
	if activity == "" {
		return 0, "", 0, fmt.Errorf("invalid activity: empty string")
	}

	timeData, err := time.ParseDuration(slicesData[2])
	if err != nil {
		return 0, "", 0, err
	}

	if timeData <= 0 {
		return 0, "", 0, fmt.Errorf("invalid time duration: %s", timeData)
	}

	return steps, activity, timeData, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	return (float64(steps) * stepLength) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepsData, activityData, timeData, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64

	switch activityData {
	case "Бег":
		if calories, err = RunningSpentCalories(
			stepsData,
			weight,
			height,
			timeData,
		); err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		if calories, err = WalkingSpentCalories(
			stepsData,
			weight,
			height,
			timeData,
		); err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activityData,
		timeData.Hours(),
		distance(stepsData, height),
		meanSpeed(stepsData, height, timeData),
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid number of steps: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight: %f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("invalid height: %f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration: %d", duration)
	}

	averageSpeed := meanSpeed(steps, height, duration)

	result := (weight * averageSpeed * duration.Minutes()) / minInH

	return result, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	result := calories * walkingCaloriesCoefficient

	return result, nil
}
