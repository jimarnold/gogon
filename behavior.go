package main

func(game *Game) collide() {
  for i,a:= range game.elements {
    for j,b := range game.elements {
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

func(game *Game) win() {
  if game.player.isDead() {
    game.gameState = lost
	return
  }

  for _, e := range game.elements {
    if e != game.player && !e.isDead() {
      game.gameState = running
		return
    }
  }
  game.gameState = won
}
