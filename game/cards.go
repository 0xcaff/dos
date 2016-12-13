package dos

import (
	proto "github.com/caffinatedmonkey/dos/proto"
	"github.com/caffinatedmonkey/dos/utils"
)

var CardColors = []proto.CardColor{
	proto.CardColor_RED,
	proto.CardColor_YELLOW,
	proto.CardColor_GREEN,
	proto.CardColor_BLUE,
	proto.CardColor_BLACK,
}

// An ObservableList of Cards.
type Cards struct {
	utils.ObservableList
}

// Creates a new deck and populates it with the standard playing cards
func PlayingDeck() *Cards {
	deck := Cards{utils.NewObservableList()}
	deck.Populate()
	deck.ObservableList.Shuffle()
	return &deck
}

func EmptyCards() *Cards {
	return &Cards{utils.NewObservableList()}
}

func (cards *Cards) Pop() proto.Card {
	return cards.ObservableList.Pop().(proto.Card)
}

func (cards *Cards) Push(c proto.Card) {
	cards.ObservableList.Push(c)
}

// Add the standard cards to the deck.
func (cards *Cards) Populate() {
	for i := 0; i < 2; i++ {
		for _, color := range CardColors {
			// Insert Cards 1-9
			for k := int32(1); k < int32(10); k++ {
				cards.ObservableList.Push(proto.Card{
					Number: k,
					Color:  color,
					Type:   proto.CardType_NORMAL,
				})
			}

			// Insert Special Cards (Skip, Reverse, DoubleDraw)
			cards.ObservableList.Push(proto.Card{
				Number: -1,
				Color:  color,
				Type:   proto.CardType_SKIP,
			}, proto.Card{
				Number: -1,
				Color:  color,
				Type:   proto.CardType_REVERSE,
			}, proto.Card{
				Number: -1,
				Color:  color,
				Type:   proto.CardType_DOUBLEDRAW,
			})
			if i == 0 {
				cards.ObservableList.Push(proto.Card{
					Number: -1,
					Color:  proto.CardColor_BLACK,
					Type:   proto.CardType_QUADDRAW,
				}, proto.Card{
					Number: 0,
					Color:  color,
					Type:   proto.CardType_NORMAL,
				})
			} else if i == 1 {
				cards.ObservableList.Push(proto.Card{
					Number: -1,
					Color:  proto.CardColor_BLACK,
					Type:   proto.CardType_WILD,
				})
			}
		}
	}
}
