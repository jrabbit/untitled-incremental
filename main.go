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

type processor func(time.Duration) int
type Brick struct {
	Processor   processor
	OutputField *int
}

func threadExperiment(dTime time.Duration) {
	// can we pass functions, can we call them sequentially?
	//identity function
	idfn := func(t time.Duration) int { return int(t.Seconds()) }
	b := Brick{Processor: idfn, OutputField: &scoreboard.Teeth}
	bricks := []Brick{b}
	for _, brick := range bricks {
		out := brick.Processor(dTime)
		*brick.OutputField = out
	}
}

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

func blit(_ js.Value, _ []js.Value) interface{} {
	jsTime, err := nullableString(js.Global().Get("localStorage").Call("getItem", "start_time"))
	var setTime time.Time
	if err != nil {
		js.Global().Get("localStorage").Call("setItem", "start_time", time.Now().Unix())
		setTime = time.Now()
	} else {
		i, _ := strconv.ParseInt(jsTime, 10, 64)
		setTime = time.Unix(i, 0)
	}
	dTime := time.Now().Sub(setTime)
	log.Println(dTime)
	log.Println(setTime)
	threadExperiment(dTime)
	scoreboard.Teeth = 1 * int(dTime.Seconds())
	query("nav > h2.teeth").Set("textContent", fmt.Sprintf("%v teeth", scoreboard.Teeth))
	fmt.Printf("%v hats", scoreboard.Hats)
	return nil
}

func secondRun() {
	const UpdateFreq = 1000
	query("#game-area").Call("removeChild", query(".game-init"))
	cb := js.FuncOf(blit)
	done := make(chan bool)
	js.Global().Get("window").Call("setInterval", cb, UpdateFreq)
	<-done
}

// PerformanceObserver or https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/setInterval

func main() {
	scoreboard = Scoreboard{}
	jsTime := js.Global().Get("localStorage").Call("getItem", "start_time")
	if jsTime == js.Null() {
		//user's first time
		firstRun()
	} else {
		secondRun()
	}
}
