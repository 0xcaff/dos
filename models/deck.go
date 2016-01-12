package dos

type Deck struct {
	Cards ObservableList
}

// Creates a new deck and populates it with the standard playing cards
func NewDeck() Deck {
	d := Deck{}

	for i := 0; i < 2; i++ {
		for _, color := range CardColors {
			// Insert Cards 1-9
			for k := 1; k < 10; k++ {
				d.Cards.Push(Card{
					Number:    k,
					CardColor: color,
					CardType:  Normal,
				})
			}

			// Insert Special Cards (Skip, Reverse, DoubleDraw)
			d.Cards.Push(Card{
				Number:    -1,
				CardColor: color,
				CardType:  Skip,
			}, Card{
				Number:    -1,
				CardColor: color,
				CardType:  Reverse,
			}, Card{
				Number:    -1,
				CardColor: color,
				CardType:  DoubleDraw,
			})
			if i == 0 {
				d.Cards.Push(Card{
					Number:    -1,
					CardColor: Black,
					CardType:  QuadDraw,
				}, Card{
					Number:    0,
					CardColor: color,
					CardType:  Normal,
				})
			} else if i == 1 {
				d.Cards.Push(Card{
					Number:    -1,
					CardColor: Black,
					CardType:  Wild,
				})
			}
		}
	}

	d.Cards.Shuffle()
	return d
}

func (d *Deck) Pop() Card {
	return d.Cards.Pop().(Card)
}

func (d *Deck) Push(c Card) {
	d.Cards.Push(c)
}
