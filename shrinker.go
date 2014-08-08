package main

type Shrinker struct {
	Thing
}

func NewShrinker(location, direction Vector2) *Shrinker {
	return &Shrinker{Thing{location: location, direction: direction, size: 6, targetSize: 6, color: Color4f{0, 1, 1, 1}}}
}

type Pickup struct {
	Thing
	timer      *IntervalTimer
	colors     [2]Color4f
	colorIndex int
}

var White Color4f = Color4f{1, 1, 1, 1}
var Red Color4f = Color4f{0, 1, 0, 1}

func NewPickup(location, direction Vector2) *Pickup {
	return &Pickup{Thing{location: location, direction: direction, size: 6, targetSize: 6, color: Color4f{0, 1, 0, 1}}, NewIntervalTimer(0.25), [2]Color4f{White, Red}, 0}
}

func (this *Pickup) update(elapsed float64) {
	this.Thing.update(elapsed)
	this.timer.Update(elapsed)
	if this.timer.Elapsed() {
		this.colorIndex = (this.colorIndex + 1) % 2
	}
	this.color = this.colors[this.colorIndex]
}
