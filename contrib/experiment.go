package main

import (
	"encoding/json"
	"fmt"
)
type Level3 struct {
	Name   string
	Number int
}
type Scoarboard1 struct {
	Teeth    int
	Hats     int
	BadStuff int
	Levels   []Level3
}
type Scoarboard2 struct {
	Teeth    int
	Hats     int
	BadStuff int
	Rings	int
	Levels   []Level3
}

func TestAddMember (){
	s1 := Scoarboard1{1,2,3, []Level3{Level3{}}}
	s, err := json.Marshal(s1)
	if err != nil {
		panic(err)
	}
	y := Scoarboard2{}
	json.Unmarshal(s, &y)
	fmt.Printf("%+v", y)
}
func main(){
	TestAddMember()
}