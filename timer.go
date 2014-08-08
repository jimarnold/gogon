package main

type IntervalTimer struct {
	interval   float64
	lastUpdate float64
	elapsed    bool
}

func NewIntervalTimer(interval float64) *IntervalTimer {
	return &IntervalTimer{interval: interval}
}

func (this *IntervalTimer) Elapsed() bool {
	return this.elapsed
}

func (this *IntervalTimer) Update(t float64) {
	this.lastUpdate += t

	if this.lastUpdate > this.interval {
		this.lastUpdate = 0
		this.elapsed = true
	} else {
		this.elapsed = false
	}
}

func (this *IntervalTimer) Reset() {
	this.lastUpdate = this.interval
}
