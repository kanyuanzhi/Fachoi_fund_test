package util

import "time"

func StartTimerByDay(f func(), hour int, min int, sec int) {
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, min, sec, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		f()
	}
}

func StartTimerByInterval(f func(), interval int) {
	for {
		t := time.NewTimer(time.Duration(interval) * time.Second)
		<-t.C
		f()
	}
}
