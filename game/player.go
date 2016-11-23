package dos

import (
	"encoding/json"
)

type Player struct {
	Name string
	Hand ObservableList
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
		Hand: NewObservableList(),
	}
}
