package main

import "fmt"
import "os"

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
type GameCondition int

func win() {
  if player.isDead() {
    fmt.Println("LOSE!")
    os.Exit(1)
  }

  for _, e := range elements {
    if e != player && !e.isDead() {
      return
    }
  }

  fmt.Println("WIN!")
  os.Exit(0)
}
