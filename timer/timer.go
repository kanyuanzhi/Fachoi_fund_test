package timer

import (
	"Fachoi_fund_test/util"
	"time"
)

type Timer struct {
	funcChan   chan *func()
	funcToInfo map[*func()]*FuncInfo
}

func NewTimer(num int) *Timer {
	// num表示最多可添加的任务数量
	return &Timer{
		make(chan *func(), num),
		make(map[*func()]*FuncInfo),
	}
}

// 添加每日运行的任务
func (t *Timer) AddDayJob(f func(), hour int, min int, sec int) {
	t.funcChan <- &f
	funcInfo := &FuncInfo{
		"day",
		hour,
		min,
		sec,
		0,
		0,
	}
	t.funcToInfo[&f] = funcInfo
}

// 添加每月固定日运行的任务
func (t *Timer) AddMonthJob(f func(), date int) {
	t.funcChan <- &f
	funcInfo := &FuncInfo{
		"month",
		0,
		0,
		0,
		0,
		date,
	}
	t.funcToInfo[&f] = funcInfo
}

// 添加间隔一段时间运行的任务
func (t *Timer) AddIntervalJob(f func(), interval int) {
	t.funcChan <- &f
	funcInfo := &FuncInfo{
		"interval",
		0,
		0,
		0,
		interval,
		0,
	}
	t.funcToInfo[&f] = funcInfo
}

func (t *Timer) Run() {
	for f := range t.funcChan {
		switch t.funcToInfo[f].jobType {
		case "day":
			util.StartTimerByDay(*f, t.funcToInfo[f].hour, t.funcToInfo[f].min, t.funcToInfo[f].sec)
		case "month":
			util.StartTimerByMonth(*f, t.funcToInfo[f].date)
		case "interval":
			util.StartTimerByInterval(*f, t.funcToInfo[f].interval)
		default:
			continue
		}
	}
	close(t.funcChan)
	for {
		time.Sleep(time.Hour * 24)
	}
}

type FuncInfo struct {
	jobType  string
	hour     int
	min      int
	sec      int
	interval int //单位秒
	date     int //单位日
}
