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

func TestStartingDeck(t *testing.T) {
	deck := PlayingDeck()

	if len(deck.List) != 108 {
		t.Errorf("Length (%d) != 108", len(deck.List))
	}

	redCards := 0
	yellowCards := 0
	greenCards := 0
	blueCards := 0
	blackCards := 0

	normalCards := 0
	skipCards := 0
	reverseCards := 0
	drawTwo := 0
	drawFour := 0
	wild := 0

	for _, card := range deck.List {
		if card.Color == proto.CardColor_RED {
			redCards++
		} else if card.Color == proto.CardColor_YELLOW {
			yellowCards++
		} else if card.Color == proto.CardColor_GREEN {
			greenCards++
		} else if card.Color == proto.CardColor_BLUE {
			blueCards++
		} else if card.Color == proto.CardColor_BLACK {
			blackCards++
		}

		if card.Type == proto.CardType_NORMAL {
			normalCards++
		} else if card.Type == proto.CardType_SKIP {
			skipCards++
		} else if card.Type == proto.CardType_REVERSE {
			reverseCards++
		} else if card.Type == proto.CardType_DOUBLEDRAW {
			drawTwo++
		} else if card.Type == proto.CardType_QUADDRAW {
			if card.Color != proto.CardColor_BLACK {
				t.Error("QuadDraw not black")
			}

			drawFour++
		} else if card.Type == proto.CardType_WILD {
			if card.Color != proto.CardColor_BLACK {
				t.Error("Wildcard not black")
			}

			wild++
		}
	}

	if redCards != 25 || yellowCards != 25 || greenCards != 25 || blueCards != 25 {
		t.Error("Wrong number of cards by color")
	}

	if blackCards != 8 {
		t.Error("Wrong number of black cards")
	}

	if wild != 4 || drawFour != 4 || drawTwo != 8 || reverseCards != 8 || skipCards != 8 {
		t.Error("Incorrect number of special cards")
	}

	if normalCards != 76 {
		t.Error("Incorrect number of normal cards")
	}
}

func TestInvalidCard(t *testing.T) {
	deck := PlayingDeck()

	card, index := deck.FindById(int32(1000))
	if card != nil || index != -1 {
		t.Error("FindById found invalid id")
	}
}

func TestEmptyPopFront(t *testing.T) {
	cards := NewCardCollection()
	card := cards.PopFront(1)
	if card != nil {
		t.Error("Empty deck popping should be nil", card)
	}
}

func TestEmptyPop(t *testing.T) {
	cards := NewCardCollection()
	card := cards.PopN(1)
	if card != nil {
		t.Error("Empty deck popping should be nil", card)
	}
}
