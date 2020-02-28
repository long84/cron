package cron

import (
	"log"
	"sync"
	"time"
)

var timers []Timer
var stop chan interface{}
var wg sync.WaitGroup

func init() {
	stop = make(chan interface{})
}

func Add(intv time.Duration, fn func()) {
	timers = append(timers, newIntervalTimer(intv, fn))
}

func AddMinutely(fn func()) {
	Add(time.Minute, fn)
}

func AddHourly(fn func()) {
	Add(time.Hour, fn)
}

func AddDaily(hour, min int, fn func()) {
	timers = append(timers, newDailyTimer(hour, min, fn))
}

func Start() {
	for _, t := range timers {
		wg.Add(1)
		go func(t Timer) {
			for {
				select {
				case <-stop:
					wg.Done()
					return
				case <-t.C():
					t.Do()
					t.Next()
				}
			}
		}(t)
	}
	log.Printf("Started %v timers", len(timers))
}

func Stop() {
	close(stop)
	wg.Wait()
	log.Println("Stopped all timers")
}