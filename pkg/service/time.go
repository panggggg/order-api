package service

import (
	"time"
)

var timerInstance *timer
var now = time.Now

type Time interface {
	Now() time.Time
	Sleep(d time.Duration)
	Freeze()
	Unfreeze()
}

type timer struct {
}

func NewTime() *timer {
	if timerInstance != nil {
		return timerInstance
	}
	timerInstance = &timer{}
	return timerInstance
}

func (t timer) Now() time.Time {
	return now()
}

func (t timer) DateFormat(datetime time.Time) string {
	return datetime.Format("2006-01-02")
}

func (t timer) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (t timer) Freeze() {
	n := time.Now()
	now = func() time.Time {
		return n
	}
}

func (t timer) Unfreeze() {
	now = time.Now
}

func (t timer) SetNow(date string) {
	n, _ := time.Parse("2006-01-02 03:04:05", date)
	now = func() time.Time {
		return n
	}
}
