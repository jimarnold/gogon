package main

import (
  glfw "github.com/go-gl/glfw3"
  "github.com/jimarnold/gltext"
)

type Profiler struct {
  frameSpeed float64
  frameTimes []float64
  frame int
  lastTime float64
  font *gltext.Font
}

func NewProfiler(font *gltext.Font) *Profiler {
  return &Profiler{frameTimes : make([]float64, 256), frameSpeed : 0, frame : 0, font: font}
}

func(this *Profiler) start() {
  this.lastTime = glfw.GetTime()
}

func(this *Profiler) update() {
  frameCount := len(this.frameTimes)
  now := glfw.GetTime()
  this.frameTimes[this.frame % frameCount] = now - this.lastTime
  movingAverage := 0.0
  for _,f := range this.frameTimes {
    movingAverage += f
  }
  movingAverage /= float64(frameCount)
  if this.frame % frameCount == 0 {
    this.frameSpeed = movingAverage
  }
  this.lastTime = now
  this.frame++
}

func(this *Profiler) render() {
  this.font.Printf(1.7, 0,"%f", this.frameSpeed)
}
