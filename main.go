package main

import (
	"time"
	"syscall/js"
	"fmt"
)

type Level struct {
	Name string
	Number int
}

type Scoreboard struct {
	Teeth float64
	Hats int
	Levels []Level
}

var scoreboard Scoreboard

type processor func(time.Duration)


func threadexperiment(d_time time.Duration){
	// can we pass functions, can we call them sequentially?
	//identity function
	idfn := func(t time.Duration){scoreboard.Teeth = t.Seconds()}
	var funcs []processor
	funcs[0] = idfn
	for _, f := range funcs {
		f(d_time)
	}
}

func firstRun(){
	js.Global().Get("localStorage").Call("setItem", "start_time", time.Now().Unix())
}


func main() {
	scoreboard = Scoreboard{}
	js_time := js.Global().Get("localStorage").Call("getItem", "start_time")
	if js_time == js.Null(){
		//user's first time
		firstRun()
	}
	//set_time, _ := time.Parse(time.RFC3339, js_time)
	fmt.Println(js_time)
	//d_time := time.Now().Sub(set_time)
	//fmt.Printf("time since game init: %d sec", d_time)
	// g, _ := time.Parse("2009", "2006")
	//threadexperiment(time.Now().Sub(g))
}