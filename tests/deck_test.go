package test

import (
	"github.com/caffinatedmonkey/dos/models"
	"testing"
)

func TestDeck(t *testing.T) {
	d := dos.NewDeck()
	if d.Cards.Length() != 108 {
		t.Errorf("Cards Length (%d) != 108", d.Cards.Length())
	}

	// blk := 0
	// red := 0
	// for _, c := range d.Cards {
	// 	if c.CardColor == dos.Black {
	// 		blk++
	// 	}

	// 	if c.CardColor == dos.Red {
	// 		red++
	// 	}
	// }

	// if blk > 8 {
	// 	t.Errorf("Black Cards (%d) > 8", blk)
	// }

	// if red > 25 {
	// 	t.Errorf("Red Cards (%d) > 25", red)
	// }
}
