package dos

import "math/rand"

type Game struct {
	Players        ObservableList
	LastPlayedCard Card
	Deck

	currentPlayerIndex int
	isReversed         bool
}

// Creates a new game an initalizes its values
func NewGame() Game {
	g := Game{}
	g.Deck = NewDeck()
	g.LastPlayedCard = g.Deck.Pop()
	g.currentPlayerIndex = -1
	return g
}

// Creates a new player, populates its deck and adds it to the game
func (g *Game) NewPlayer(name string) Player {
	p := Player{Name: name}
	for i := 0; i < 8; i++ {
		c := g.Deck.Pop()
		p.Hand.Push(c)
	}
	g.Players.Push(p)

	return p
}

// Called after a player completes their turn. Get's the player who is to play
// next.
func (g *Game) NextPlayer() Player {

	// TODO: More Cards
	increment := 1
	switch g.LastPlayedCard.CardType {
	case Reverse:
		g.isReversed = !g.isReversed

	case Skip:
		increment = increment + 1
	}

	// If this is the first turn, pick a random player to start.
	if g.currentPlayerIndex == -1 {
		g.currentPlayerIndex = rand.Intn(g.Players.Length())
	}

	if g.isReversed {
		increment = increment * -1
	} else {
		increment = increment * 1
	}

	// Causes the players to be run in a cycle
	index := increment + g.currentPlayerIndex
	if index < 0 {
		index = g.Players.Length() + index
	} else if index >= g.Players.Length() {
		index = g.Players.Length() - index
	}
	g.currentPlayerIndex = index
	return g.Players.GetElement(index).(Player)
}

func (g *Game) PlayCard(c Card) {
	if g.LastPlayedCard.CanCover(c) {
		g.LastPlayedCard = c
	}
}
