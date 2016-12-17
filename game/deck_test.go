package dos

// import (
// 	"github.com/caffinatedmonkey/dos/models"
// 	"testing"
// )
//
// var deckCounts map[dos.CardColor]int = map[dos.CardColor]int{
// 	dos.Blue:   25,
// 	dos.Green:  25,
// 	dos.Red:    25,
// 	dos.Orange: 25,
// 	dos.Black:  8,
// }
//
// func TestStartingDeck(t *testing.T) {
// 	d := dos.NewDeck()
// 	if d.Cards.Length() != 108 {
// 		t.Errorf("Length (%d) != 108", d.Cards.Length())
// 	}
//
// 	colorCount := make(map[dos.CardColor]int)
// 	d.Cards.Each(func(v interface{}, i int) bool {
// 		c := v.(dos.Card)
// 		colorCount[c.CardColor] = colorCount[c.CardColor] + 1
// 		return false
// 	})
//
// 	for k, v := range colorCount {
// 		if deckCounts[k] != v {
// 			t.Logf("%d (should be) != %d (counted)", deckCounts[k], v)
// 		}
// 	}
//
// 	// TODO: Count Numbers
// 	// TODO: Count Type
// }
