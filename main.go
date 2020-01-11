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
	Rings    int
	Levels   []Level
}

func (s *Scoreboard) BadStuffHappened(stuff int) {
	s.BadStuff = stuff + s.BadStuff
}

var scoreboard Scoreboard

const TeethPerChild = 20
const TeethPerAdult = 32

type PlanetX struct {
	Kids              int
	SonicFactor       int
	Hidden            bool
	ScreamingDecibels int
	LastCheckinTime   time.Time `json:",omitempty"`
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
var sonicModeEnabled = false

type DataItem struct {
	Key  string
	Item interface{}
}

var DataStore = []DataItem{DataItem{"planet_state", planetx}, {"sonic_state", sonic}}

func UnifiedStorageSave(d []DataItem) {
	for _, dataitem := range d {
		outJson, err := json.Marshal(dataitem)
		if err != nil {
			log.Println(err)
		}
		js.Global().Get("localStorage").Call("setItem", dataitem.Key, string(outJson))
	}
}

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
	// log.Printf("%+v scoreboard", scoreboard)
	return nil
}

type Position struct {
	x int
	y int
}

type Sonic struct {
	Position Position
}

var sonic Sonic

func keyDown(_ js.Value, args []js.Value) interface{} {
	event := args[0]
	switch key := event.Get("key").String(); key {
	case "ArrowLeft":
		sonic.Position.x = -1
	case "ArrowRight":
		sonic.Position.x += 1
	case "ArrowUp":
		sonic.Position.y += 1
	case "ArrowDown":
		sonic.Position.y -= 1
	default:
		log.Printf("random fucking key? %v", key)
	}
	return nil
}

func sonicTime() {
	sect := query("section.sonic")
	sect.Call("removeAttribute", "hidden")
	// keybinds
	cb := js.FuncOf(keyDown)
	js.Global().Get("document").Call("addEventListener", "keydown", cb)

}

func askSonic() {
	button := js.Global().Get("document").Call("createElement", "button")
	button.Set("textContent", "Sonic?")
	query("#game-area").Call("append", button)
	cb := js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		sonicModeEnabled = true
		return nil
	})
	button.Call("addEventListener", "click", cb)
}

func secondRun() {
	const UpdateFreq = 1000
	query("#game-area").Call("removeChild", query(".game-init"))
	cb := js.FuncOf(blit)
	planetXCB := js.FuncOf(procPlanetX)
	done := make(chan bool)
	js.Global().Get("window").Call("setInterval", cb, UpdateFreq)
	js.Global().Get("window").Call("setInterval", planetXCB, 5*UpdateFreq)
	askSonic()
	if sonicModeEnabled {
		sonicTime()
	}
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
