package main
func collide() {
  for i,a:= range elements {
    for j,b := range elements {
      if j <= i {
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
