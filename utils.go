package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

var LOGGING = flag.Bool("l", false, "log to console")

func init() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
}

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

type RGB struct{ r, g, b float64 }

func clamp(f float64, min, max float64) float64 {
	if f > max {
		return max
	} else if f < min {
		return min
	}
	return f
}

type Rect struct {
	top, bottom, left, right float64
}

func debug(v ...interface{}) {
	if *LOGGING {
		log.Println(v)
	}
}

func debugf(s string, v ...interface{}) {
	if *LOGGING {
		log.Printf(s, v)
	}
}
