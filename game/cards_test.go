package dos

import (
	proto "github.com/caffinatedmonkey/dos/proto"
	"testing"
)

func TestRemoveCard(t *testing.T) {
	cards := NewCardCollection()
	for i := 0; i < 9; i++ {
		cards.Push(proto.Card{
			Number: int32(i),
			Type:   proto.CardType_NORMAL,
			Color:  proto.CardColor_RED,
		})
	}

	cards.RemoveCard(5)

	for _, card := range cards.List {
		if card.Number == int32(5) {
			t.Fail()
		}
	}
}
