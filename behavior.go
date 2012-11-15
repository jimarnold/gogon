package main

func collide() {
  for _,a:= range elements {
    for _,b := range elements {
      if a == b || (a.isDead() || b.isDead()) {
        continue
      }

      if a.intersects(b) {
        if a.biggerThan(b) {
          a.absorb(b)
        } else {
          b.absorb(a)
        }
      }
    }
  }
}

func win() GameState {
  if player.isDead() {
    return lost
  }

  for _, e := range elements {
    if e != player && !e.isDead() {
      return running
    }
  }

  return won
}
