package main

import (
	"log"
	"syscall/js"
)

type Position struct {
	x int
	y int
}

type Sonic struct {
	Position Position
	Snake    bool
	Asking   bool
}

func keyDown(_ js.Value, args []js.Value) interface{} {
	event := args[0]
	switch key := event.Get("key").String(); key {
	case "ArrowLeft":
		sonic.Position.x = -1
	case "ArrowRight":
		sonic.Position.x += 1
	case "ArrowUp":
		if sonic.Snake {
			sonic.Position.y += 1
		}
	case "ArrowDown":
		if sonic.Snake {
			sonic.Position.y -= 1
		}
	default:
		log.Printf("random fucking key? %v", key)
	}
	return nil
}

func sonicTime() {
	//init sonic
	sect := query("section.sonic")
	sect.Call("removeAttribute", "hidden")
	// key binds done here
	cb := js.FuncOf(keyDown)
	js.Global().Get("document").Call("addEventListener", "keydown", cb)
}

func askSonic() bool {
	// returns if added element.
	if sonicModeEnabled {
		return false
	}
	if sonic.Asking {
		return false
	}
	button := js.Global().Get("document").Call("createElement", "button")
	button.Set("textContent", "Sonic?")
	query("#game-area").Call("append", button)
	cb := js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		sonicModeEnabled = true
		sonicTime()
		return nil
	})
	button.Call("addEventListener", "click", cb)
	sonic.Asking = true
	return true
}
