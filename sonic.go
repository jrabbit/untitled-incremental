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
}

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
