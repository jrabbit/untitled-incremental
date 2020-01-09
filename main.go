package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"syscall/js"
	"time"
)

type Level struct {
	Name   string
	Number int
}

type Scoreboard struct {
	Teeth    int
	Hats     int
	BadStuff int
	Levels   []Level
}

func (s *Scoreboard) BadStuffHappened(stuff int) {
	s.BadStuff = stuff + s.BadStuff
}

var scoreboard Scoreboard

const TeethPerChild = 20
const TeethPerAdult = 32

type PlanetX struct {
	Kids            int
	SonicFactor     int
	Hidden          bool
	LastCheckinTime time.Time `json:",omitempty"`
}

func (p *PlanetX) load() {
	jsState, err := nullableString(js.Global().Get("localStorage").Call("getItem", "planet_state"))
	if err != nil {
		log.Print(err)
	}
	jsonErr := json.Unmarshal([]byte(jsState), p)
	if jsonErr != nil {
		log.Print(jsonErr)
	}
}

func (p *PlanetX) CheckIn() {
	//previousTime := p.LastCheckinTime
	now := time.Now()
	//dT := previousTime.Sub(now)
	// do population model here
	newKids := math.Pow(float64(p.Kids), 0.26)
	p.Kids = p.Kids + int(newKids)
	p.LastCheckinTime = now
	p.save()
}

func (p *PlanetX) getTeeth(number int) bool {
	currTeeth := p.Kids * TeethPerChild
	if currTeeth > number {
		targetTeeth := currTeeth - number
		p.Kids = targetTeeth / TeethPerChild
		badStuff := targetTeeth % TeethPerChild
		scoreboard.BadStuffHappened(badStuff)
		return true
	} else {
		return false
	}
}
func (p *PlanetX) save() {
	outJson, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	} else {
		js.Global().Get("localStorage").Call("setItem", "planet_state", string(outJson))
	}
}

var planetx PlanetX

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
		planetx.Kids = 500
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

func procPlanetX(_ js.Value, _ []js.Value) interface{} {
	planetx.CheckIn()
	log.Printf("planetx: %+v", planetx)
	return nil
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
	//log.Printf("%d dTime", dTime)
	//log.Printf("%v game started", setTime)
	threadExperiment(dTime)
	query("nav > h2.teeth").Set("textContent", fmt.Sprintf("%v teeth", scoreboard.Teeth))
	log.Printf("%+v scoreboard", scoreboard)
	return nil
}

func secondRun() {
	const UpdateFreq = 1000
	query("#game-area").Call("removeChild", query(".game-init"))
	cb := js.FuncOf(blit)
	planetXCB := js.FuncOf(procPlanetX)
	done := make(chan bool)
	js.Global().Get("window").Call("setInterval", cb, UpdateFreq)
	js.Global().Get("window").Call("setInterval", planetXCB, 5*UpdateFreq)
	<-done
}

// PerformanceObserver or https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/setInterval

func main() {
	scoreboard = Scoreboard{}
	planetx = PlanetX{Kids: 500, Hidden: true, SonicFactor: 1}
	jsTime := js.Global().Get("localStorage").Call("getItem", "start_time")
	if jsTime == js.Null() {
		//user's first time
		firstRun()
	} else {
		planetx.load()
		secondRun()
	}
}
