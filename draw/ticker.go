package draw

import (
	"time"
)

// Ticker calls targets as close to it's interval as possible, but will not call a target until its
// previous call has returned.
type Ticker struct {
	target func(elapsed float32)
	stop   chan time.Time
	start  chan time.Time
	done   chan struct{}
}

// NewTicker creates a new Ticker with the provided target and interval. The ticker is paused by default,
// and must be started by calling Start.
func NewTicker(interval time.Duration, target func(elapsed float32)) *Ticker {
	t := &Ticker{
		stop:   make(chan time.Time),
		start:  make(chan time.Time),
		done:   make(chan struct{}),
		target: target,
	}

	go func() {
		timer := time.NewTimer(interval)
		defer timer.Stop()
		startTime := time.Now()
		lastTime := startTime

		for {
			select {
			case <-t.done:
				return
			case sTime := <-t.start:
				if sTime.After(startTime) {
					startTime = sTime
				}
			case pauseTime := <-t.stop: // Paused
				if !timer.Stop() {
					<-timer.C
				}
				startTime = t.pause(pauseTime)
				timer.Reset(interval - pauseTime.Sub(lastTime))
				lastTime = startTime
			case <-timer.C:
				currentTime := time.Now()
				elapsed := time.Now().Sub(lastTime)
				target(float32(elapsed) / float32(time.Second))

				// Recalculate elapsed to account for tick taking a long time to return
				elapsed = time.Now().Sub(lastTime)
				timer.Reset(interval - elapsed)
				lastTime = currentTime
			}
		}
	}()

	return t
}

// puase blocks until a start signal is recieved that is later than the last stop signal, and returns the
// start signal time that ended the pause. The pauseTime parameter is used as the initial pauseTime until a later
// pause signal is recieved.
func (t *Ticker) pause(pauseTime time.Time) time.Time {
	for {
		select {
		case pTime := <-t.stop:
			if pTime.After(pauseTime) {
				pauseTime = pTime
			}
		case startTime := <-t.start:
			if !startTime.Before(pauseTime) {
				return startTime
			}
		}
	}
}

// Start starts the ticker.
func (t *Ticker) Start() {
	t.start <- time.Now()
}

// Stop stops the ticker. No more ticks will be generated until Start is called.
func (t *Ticker) Stop() {
	t.stop <- time.Now()
}

// Close permenantly stops the ticker, cleaning up it's goroutine.
// Always call Close when done with a ticker so it can be garbage collected.
func (t *Ticker) Close() {
	t.start <- time.Now()
	t.done <- struct{}{}
}
