package main

import (
  "github.com/go-gl/glfw"
  //"strconv"
	//"fmt"
)

type Profiler struct {
  frameSpeed float64
  frameTimes []float64
  frame int
  startTime float64
}

func NewProfiler() *Profiler {
  return &Profiler{frameTimes : make([]float64, 32), frameSpeed : 0, frame : 0}
}

func(this *Profiler) start() {
  this.startTime = glfw.Time()
}

func(this *Profiler) stop() {
  frameCount := len(this.frameTimes)
  this.frameTimes[this.frame % frameCount] = glfw.Time() - this.startTime
  movingAverage := 0.0
  for _,f := range this.frameTimes {
    movingAverage += f
  }
  movingAverage /= float64(frameCount)
  if this.frame % frameCount == 0 {
    this.frameSpeed = movingAverage
  }

  this.frame++
}
var i int = 0
func(this *Profiler) render() {
	//i++
	//if i % 32 == 0 {
		//fmt.Println(strconv.FormatFloat(this.frameSpeed, 'f',10,32))
	//}
	//font.drawString(10,10,strconv.FormatFloat(this.frameSpeed, 'f',10,32))
}
