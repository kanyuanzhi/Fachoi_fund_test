package util

import (
	"fmt"
	"time"
)

// 每天固定时间运行
func StartTimerByDay(f func(), hour int, min int, sec int) {
	go func() {
		for {
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), hour, min, sec, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			f()
			fmt.Printf("StartTimerByDay: %dh%dmin%ds\n", hour, min, sec)
			fmt.Println(time.Now())
		}
	}()
}

// 每间隔固定时间运行
func StartTimerByInterval(f func(), interval int) {
	go func() {
		for {
			t := time.NewTimer(time.Duration(interval) * time.Second)
			<-t.C
			f()
			fmt.Printf("StartTimerByInterval: %ds\n", interval)
			fmt.Println(time.Now())
		}
	}()
}

// 每个月固定日期运行
func StartTimerByMonth(f func(), date int) {
	go func() {
		for {
			now := time.Now()
			next := now.AddDate(0, 1, 0)
			next = time.Date(next.Year(), next.Month(), date, 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			f()
			fmt.Printf("StartTimerByMonth: %dd\n", date)
			fmt.Println(time.Now())
		}
	}()
}
