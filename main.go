package main

import (
	"fmt"
	"syscall/js"
	"time"
)

type Level struct {
	Name   string
	Number int
}

type Scoreboard struct {
	Teeth  float64
	Hats   int
	Levels []Level
}

var scoreboard Scoreboard

type processor func(time.Duration)

//func threadexperiment(d_time time.Duration) {
//	// can we pass functions, can we call them sequentially?
//	//identity function
//	idfn := func(t time.Duration) { scoreboard.Teeth = t.Seconds() }
//	var funcs []processor
//	funcs[0] = idfn
//	for _, f := range funcs {
//		f(d_time)
//	}
//}

func firstRun() {
	cb := func() {
		js.Global().Get("localStorage").Call("setItem", "start_time", time.Now().Unix())
		fmt.Println("we started!!!")
	}
	query(".game-init").Call("addEventListener", "click", cb)
}

func query(qs string) js.Value {
	return js.Global().Get("document").Call("querySelector", qs)
}

func blit() {
	query("nav > h2.teeth").Call("textContent", fmt.Sprintf("%d teeth", scoreboard.Teeth))
}

func secondRun(d_time time.Duration) {
	query("#game-area").Call("removeChild", query(".game-init"))
}

func main() {
	scoreboard = Scoreboard{}
	js_time := js.Global().Get("localStorage").Call("getItem", "start_time")
	if js_time == js.Null() {
		//user's first time
		firstRun()
	} else {
		set_time, _ := time.Parse(time.UnixDate, js_time.String())
		d_time := time.Now().Sub(set_time)
		secondRun(d_time)
		blit()
	}
	fmt.Println(js_time)
	//fmt.Printf("time since game init: %d sec", d_time)
	//threadexperiment(time.Now().Sub(g))
}
