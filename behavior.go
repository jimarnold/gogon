package main

import "fmt"
import "os"

func collide() {
  for _,a:= range things {
    for _,b := range things {
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
type GameCondition int

func win() {
  if player.isDead() {
    fmt.Println("LOSE!")
    os.Exit(1)
  }

  for _, thing := range things {
    if thing != player && !thing.isDead() {
      return
    }
  }

  fmt.Println("WIN!")
  os.Exit(0)
}
