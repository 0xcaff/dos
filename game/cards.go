package dos

type Card struct {
	Number int
	CardType
	CardColor
}

type CardType string

const (
	Normal     CardType = "normal"
	Skip                = "skip"
	DoubleDraw          = "drawtwo"
	Reverse             = "reverse"
	Wild                = "wild"
	QuadDraw            = "wilddraw"
)

type CardColor string

const (
	Red    CardColor = "red"
	Orange           = "orange"
	Green            = "green"
	Blue             = "blue"
	Black            = "black"
)

var CardColors []CardColor = []CardColor{Red, Orange, Green, Blue}

// Returns whether or not c can be played on top of currentCard
func (baseCard *Card) CanCover(otherCard Card) bool {
	if baseCard.CardType != Normal {
		specialMatch := baseCard.CardType == otherCard.CardType
		if specialMatch {
			// Matching Special Cards
			return true
		}
	}

	colorsMatch := baseCard.CardColor == otherCard.CardColor
	numbersMatch := baseCard.Number == otherCard.Number

	// Black Cards can be played at any time.
	return colorsMatch || numbersMatch || otherCard.CardColor == Black
}
