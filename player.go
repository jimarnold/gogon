package main

type Player struct {
  Thing
}

const speed = 0.1
const maxSpeed float64 = 2.0
const brake = .95

var keyMap map[int] KeyAction

func init() {
  keyMap = map[int] KeyAction {
    KeyW : up,
    KeyS : down,
    KeyA : left,
    KeyD : right,
  }
}

func (this *Player) update(elapsed float64) {
  this.Thing.update(elapsed)

  handleKeys(keyMap)

  this.direction = this.direction.scale(brake)
}

func up() {
  player.thrust(Vector2{0,-speed})
}

func down() {
  player.thrust(Vector2{0,speed})
}

func left() {
  player.thrust(Vector2{-speed,0})
}

func right() {
  player.thrust(Vector2{speed,0})
}

func(this *Player) thrust(v Vector2) {
  this.direction = this.direction.Add(v).clampedTo(maxSpeed)
}

