package main

import (
	"fmt"
	"math"
	"time"
)

var (
	now     = time.Now()
	h1Start = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	h2Start = time.Date(now.Year(), time.July, 1, 0, 0, 0, 0, time.Local)
)

func hStart() (string, time.Time) {
	if now.After(h2Start) {
		return fmt.Sprintf("2H%d", now.Year()), h2Start
	}
	return fmt.Sprintf("1H%d", now.Year()), h1Start

}

func main() {
	half, start := hStart()
	end := start.AddDate(0, 6, 0)
	fromStart := now.Sub(start)
	tilEnd := -now.Sub(end)

	fmt.Printf("%s Week %d: weeks left: %d, days left: %d, approx work days left: %d \n",
		half, weekOf(fromStart), weeks(tilEnd), days(tilEnd), workdays(tilEnd))
}

func days(d time.Duration) int {
	return int(d.Hours()) / 24
}

func workdays(d time.Duration) int {
	return weeks(d) * 5
}

func weeks(d time.Duration) int {
	return days(d) / 7.0
}

func weekOf(d time.Duration) int {
	return int(math.Ceil(d.Hours() / 24 / 7.0))
}
