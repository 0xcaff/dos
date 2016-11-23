package dos

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"github.com/gorilla/websocket"
)

type Game struct {
	Players        ObservableList
	LastPlayedCard Card
	Deck

	currentPlayerIndex int
	isReversed         bool
	lastPlayerPlayed   bool
}

// Creates a new game an initalizes its values
func NewGame() Game {
	// Initalize Values
	g := Game{
		Players:            NewObservableList(),
		Deck:               NewDeck(),
		currentPlayerIndex: -1,
		lastPlayerPlayed:   true,
	}
	g.LastPlayedCard = g.Deck.Pop()
	return g
}

// Creates a new player, populates its deck and adds it to the game
func (g *Game) NewPlayer(name string) *Player {
	p := NewPlayer(name, conn)
	g.DrawCards(p, 8)
	g.Players.Push(p)
	return p
}

// Called after a player completes their turn. Get's the player who is to play
// next.
func (g *Game) NextPlayer() *Player {
	// If this is the first turn, pick a random player to start.
	if g.currentPlayerIndex == -1 {
		g.currentPlayerIndex = rand.Intn(g.Players.Length())
	}

	increment := 1
	if g.lastPlayerPlayed {
		switch g.LastPlayedCard.CardType {
		case Reverse:
			g.isReversed = !g.isReversed
			if g.Players.Length() == 2 {
				increment = increment + 1
			}

		case Skip:
			increment = increment + 1

		case DoubleDraw:
			p, _ := g.GetPlayer(g.currentPlayerIndex + increment)
			g.DrawCards(p, 2)
			increment = increment + 1

		case QuadDraw:
			p, _ := g.GetPlayer(g.currentPlayerIndex + increment)
			g.DrawCards(p, 2)
			increment = increment + 1
		}
	}

	p, cycint := g.GetPlayer(g.currentPlayerIndex + increment)
	g.currentPlayerIndex = cycint
	return p
}

func (g *Game) PlayCard(p *Player, index int) bool {
	c := p.Hand.GetElement(index).(Card)
	if g.LastPlayedCard.CanCover(c) {
		p.Hand.RemoveElement(index)
		g.LastPlayedCard = c
		return true
	} else {
		return false
	}
}

// Draws count cards into p's hand.
func (g *Game) DrawCards(p *Player, count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		c := g.Deck.Pop()
		p.Hand.Push(c)
	}
}

func (g *Game) Start() error {
	if g.Players.Length() < 2 {
		return fmt.Errorf("Not enough players %d\n", g.Players.Length())
	}

	// Send player information to other players

gameLoop:
	for {
		player := g.NextPlayer()

		err := g.SendToPlayers(map[string]interface{}{
			"type":   "update",
			"for":    "turn",
			"what":   g.LastPlayedCard,
			"active": player.Name,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[game] (%s) turn\n", player.Name)

		drewCard := false

	responseLoop:
		for {
			// Get Response From Player
			resp, err := player.getResponse()
			if err != nil {
				fmt.Printf("[websocket] (error) %v\n", err)
				break
			}

			if resp.Type == "draw" && !drewCard {
				g.DrawCards(player, 1)
				drewCard = true
			} else if resp.Type == "play" {
				fmt.Printf("[game] (%s) hand: %+v\n", player.Name, player.Hand.List)
				c := player.Hand.GetElement(resp.Card)
				r := g.PlayCard(player, resp.Card)

				if r {
					fmt.Printf("[game] (%s) played (%v)\n", player.Name, c)
					g.lastPlayerPlayed = true
					break responseLoop
				} else {
					player.Conn.WriteJSON(map[string]string{
						"type":    "error",
						"message": "You can't play that card",
					})
					fmt.Printf("[game] (%s) failed to play card (%v)\n", player.Name, c)
				}
			} else {
				fmt.Printf("[game] (%s) sent invalid command %s\n", player.Name, resp.Type)
				continue responseLoop
			}

			// If a card has been drew and there isn't a card to play, continue to the
			// next player.
			if drewCard {
				for _, v := range player.Hand.List {
					if t := v.(Card); g.LastPlayedCard.CanCover(t) {
						fmt.Printf("[game] (%s) can play (%v)\n", player.Name, v)
						continue responseLoop
					}
				}
				g.lastPlayerPlayed = false
				break responseLoop
			}
		}

		for _, v := range g.Players.List {
			p := v.(*Player)
			if p.Hand.Length() < 1 {
				fmt.Printf("[game] %s wins\n", p.Name)
				// TODO: Implement on Client
				g.SendToClients(map[string]string{
					"type":   "end",
					"winner": p.Name,
				})
				break gameLoop
			}
		}

		// TODO: This is trash.
		if g.Deck.Cards.Length() < 5 {
			g.Deck.Cards.Clear()
			g.Deck.Populate()
		}
	}
}

func (g *Game) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     "init",
		"deck":     g.Deck.Cards.List,
		"players":  g.Players.List,
		"lastCard": g.LastPlayedCard,
	})
}

func (g *Game) EachPlayer(f func(*Player) bool) bool {
	return g.Players.Each(func(i interface{}, j int) bool {
		player := i.(*Player)
		return f(player)
	})
}

// Gets player at index in a cyclic fashion.
func (g *Game) GetPlayer(index int) (*Player, int) {
	if g.isReversed {
		index = index * -1
	} else {
		index = index * 1
	}

	for !(index >= 0 && index < g.Players.Length()) {
		if index < 0 {
			index = g.Players.Length() + index
		} else if index >= g.Players.Length() {
			index = index - g.Players.Length()
		}
	}
	p := g.Players.GetElement(index).(*Player)
	return p, index
}

func (g *Game) SendToSpectators(i interface{}) error {
	for _, conn := range g.Spectators {
		err := conn.WriteJSON(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) SendToPlayers(i interface{}) error {
	var err error
	g.EachPlayer(func(p *Player) bool {
		err = p.Conn.WriteJSON(i)
		if err != nil {
			return true
		} else {
			return false
		}
	})
	return err
}

func (g *Game) SendToClients(i interface{}) error {
	err := g.SendToSpectators(i)
	if err != nil {
		return err
	}
	err = g.SendToPlayers(i)
	if err != nil {
		return err
	}
	return nil
}

func handle(g *Game, typ, fr string) func(interface{}) {
	return handleExtra(g, typ, fr, nil)
}

func handleExtra(g *Game, typ, fr string, extra map[string]interface{}) func(interface{}) {
	return func(i interface{}) {
		err := g.SendToSpectators(mergeMaps(map[string]interface{}{
			"type": typ,
			"for":  fr,
			"what": i,
		}, extra))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			r[k] = v
		}
	}
	return r
}
