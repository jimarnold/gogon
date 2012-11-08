package main

func collide() {
  for _,a:= range things {
    for _,b := range things {
      if a == b || (a.size == 0 || b.size == 0) {
        continue
      }

      if a.intersects(b) {
        if a.size > b.size {
          a.absorb(b)
        } else {
          b.absorb(a)
        }
      }
    }
  }
}
