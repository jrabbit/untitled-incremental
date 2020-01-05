package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"syscall/js"
	"time"
)

type Level struct {
	Name   string
	Number int
}

type Scoreboard struct {
	Teeth  int
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
	done := make(chan bool)
	cb := func(this js.Value, args []js.Value) interface{} {
		js.Global().Get("localStorage").Call("setItem", "start_time", time.Now().Unix())
		secondRun()
		fmt.Println("we started!!!")
		return nil
	}
	query(".game-init").Call("addEventListener", "click", js.FuncOf(cb))
	<-done
}

func query(qs string) js.Value {
	return js.Global().Get("document").Call("querySelector", qs)
}

var ErrJSNull = errors.New("expected string but got null")

func nullableString(thing js.Value) (string, error) {
	if thing == js.Null() {
		return "", fmt.Errorf("%q", ErrJSNull)
	}
	return thing.String(), nil
}

func blit(this js.Value, args []js.Value) interface{} {
	js_time, err := nullableString(js.Global().Get("localStorage").Call("getItem", "start_time"))
	var set_time time.Time
	if err != nil {
		js.Global().Get("localStorage").Call("setItem", "start_time", time.Now().Unix())
		set_time = time.Now()
	} else {
		i, _ := strconv.ParseInt(js_time, 10, 64)
		set_time = time.Unix(i, 0)
	}
	d_time := time.Now().Sub(set_time)
	log.Println(d_time)
	log.Println(set_time)
	scoreboard.Teeth = 1 * int(d_time.Seconds())
	query("nav > h2.teeth").Set("textContent", fmt.Sprintf("%v teeth", scoreboard.Teeth))
	fmt.Printf("%v hats", scoreboard.Hats)
	return nil
}

func secondRun() {
	const UPDATE_FREQ = 1000
	query("#game-area").Call("removeChild", query(".game-init"))
	cb := js.FuncOf(blit)
	done := make(chan bool)
	js.Global().Get("window").Call("setInterval", cb, UPDATE_FREQ)
	<-done
}

// PerformanceObserver or https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/setInterval

func main() {
	scoreboard = Scoreboard{}
	js_time := js.Global().Get("localStorage").Call("getItem", "start_time")
	if js_time == js.Null() {
		//user's first time
		firstRun()
	} else {
		secondRun()
	}
}
