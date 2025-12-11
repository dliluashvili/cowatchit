package helpers

import (
	"log"
	"time"
)

func StringToDate(value string) time.Time {
	parsedDate, err := time.Parse("2006-01-02", value)
	if err != nil {
		log.Println("Invalid date_of_birth format:", err)
	}
	return parsedDate
}

func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()

	if now.YearDay() < dob.YearDay() {
		age--
	}

	return age
}
