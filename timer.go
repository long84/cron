package cron

import "time"

type Timer interface {
	C() <-chan time.Time
	Do()
	Next()
}

type IntervalTimer struct {
	intv	time.Duration
	timer	*time.Timer
	fn		func()
}

func (t *IntervalTimer) Next() {
	now := time.Now()
	intv := int64(t.intv/time.Second)
	nextTs := (now.Unix()/intv + 1) * intv
	next := time.Unix(nextTs, 1e7)

	t.timer = time.NewTimer(next.Sub(now))
}

func (t *IntervalTimer) Do() {
	t.fn()
}

func (t *IntervalTimer) C() <-chan time.Time {
	return t.timer.C
}

func newIntervalTimer(intv time.Duration, fn func()) *IntervalTimer {
	t := &IntervalTimer{
		intv: intv,
		fn: fn,
	}
	t.Next()
	return t
}

type DailyTimer struct {
	hour 	int
	min		int
	timer	*time.Timer
	fn		func()
}

func (t *DailyTimer) Next() {
	now := time.Now()

	// 若当前没过则今日执行，否则明日执行
	var next time.Time
	if now.Hour() < t.hour || (now.Hour() == t.hour && now.Minute() < t.min) {
		next = time.Date(now.Year(), now.Month(), now.Day(), t.hour, t.min, 0, 1e7, time.Local)
	} else {
		next = now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), t.hour, t.min, 0, 1e7, time.Local)
	}

	t.timer = time.NewTimer(next.Sub(now))
}

func (t *DailyTimer) Do() {
	t.fn()
}

func (t *DailyTimer) C() <-chan time.Time {
	return t.timer.C
}

func newDailyTimer(hour, min int, fn func()) *DailyTimer {
	t := &DailyTimer{
		hour:  hour,
		min:   min,
		fn:    fn,
	}
	t.Next()
	return t
}