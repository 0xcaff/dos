package dos

import (
	"errors"
	// "fmt"
	// "log"
	"math/rand"
	"sync"

	proto "github.com/caffinatedmonkey/dos/proto"
	"github.com/caffinatedmonkey/dos/utils"
)

type Player struct {
	Name string
	Cards
	TurnDone chan interface{}
}

type Game struct {
	players []*Player
	Discard Cards
	Deck    Cards

	PlayerJoined utils.Broadcaster
	PlayerLeft   utils.Broadcaster
	Turn         utils.Broadcaster

	// TODO: Thread safe reading
	playerMutex sync.Mutex

	currentPlayerIndex int
	isReversed         bool
	// lastPlayerPlayed   bool
}

// Creates a new game an initalizes its values
func NewGame() *Game {
	// Initalize Values
	g := Game{
		players:      []*Player{},
		Deck:         *PlayingDeck(),
		PlayerJoined: *utils.NewBroadcaster(),
		PlayerLeft:   *utils.NewBroadcaster(),
		Turn:         *utils.NewBroadcaster(),

		currentPlayerIndex: -1,
		// lastPlayerPlayed:   true,
	}

	g.Discard.Push(g.Deck.Pop())
	return &g
}

// Creates a new player, populates its deck and adds it to the game
func (game *Game) NewPlayer(name string) (*Player, error) {
	for _, player := range game.players {
		if player.Name == name {
			return nil, errors.New("Player already exists in game")
		}
	}

	player := &Player{
		Cards:    *NewCardCollection(),
		Name:     name,
		TurnDone: make(chan interface{}),
	}
	game.DrawCards(&player.Cards, 8)

	game.playerMutex.Lock()
	game.players = append(game.players, player)
	game.playerMutex.Unlock()

	// Inform Listeners
	game.PlayerJoined.Broadcast(name)

	return player, nil
}

func (game *Game) RemovePlayer(removing *Player) {
	game.playerMutex.Lock()

	// Remove Player
	i := 0
	newPlayers := make([]*Player, len(game.players))
	for _, player := range game.players {
		if player != removing {
			i++
			newPlayers[i] = player
		}
	}

	game.players = newPlayers
	game.playerMutex.Unlock()

	// Return to discard pile
	game.Discard.PushFront(removing.Cards.List...)

	// Notify Players
	game.PlayerLeft.Broadcast(removing.Name)
}

func (game *Game) GetPlayerList() []string {
	result := make([]string, len(game.players))
	i := 0
	for _, player := range game.players {
		result[i] = player.Name
		i++
	}
	return result
}

// Called after a player completes their turn. Get's the player who is to play
// next.
func (game *Game) NextPlayer() *Player {
	// If this is the first turn, pick a random player
	if game.currentPlayerIndex == -1 {
		game.currentPlayerIndex = rand.Intn(len(game.players) - 1)
		return game.players[game.currentPlayerIndex]
	}

	increment := 1
	lastCard := game.Discard.List[len(game.Discard.List)-1]

	switch lastCard.Type {
	case proto.CardType_REVERSE:
		game.isReversed = !game.isReversed
		if len(game.Players) == 2 {
			increment += 1
		}

	case proto.CardType_SKIP:
		increment += 1

	case proto.CardType_DOUBLEDRAW:
		increment += 1
		player, _ := game.GetPlayer(1)
		game.DrawCards(&player.Cards, 2)

	case proto.CardType_QUADDRAW:
		increment += 1
		player, _ := game.GetPlayer(1)
		game.DrawCards(&player.Cards, 4)

	}

	nextPlayer, index := game.GetPlayer(increment)
	game.currentPlayerIndex = index
	return nextPlayer
}

// Gets the player index n positions away from the current player.
func (game *Game) GetPlayer(n int) (*Player, int) {
	if game.isReversed {
		n *= -1
	}

	current := (game.currentPlayerIndex + n) % len(game.players)
	return game.players[current], current
}

func (game *Game) PlayCard(player *Player, id int32, color proto.CardColor) error {
	card := player.Cards.PopId(id)
	if card == nil {
		return errors.New("Card is not owned by player")
	}

	// Check If Card Is Valid
	lastDiscard := &game.Discard.List[len(game.Discard.List)-1]
	canCover := CanCoverCard(lastDiscard, card)
	if !canCover {
		return errors.New("You can't play that card")
	}

	// Set Color If Needed
	if card.Color == proto.CardColor_BLACK {
		card.Color = color
	}

	game.Discard.Push(*card)
	return nil
}

// Draws count cards into hand. Recycle's cards if needed
func (game *Game) DrawCards(hand *Cards, count int) {
	if count > len(game.Deck.List) {
		// Not enough cards to draw. Recycle discard pile.
		recyclable := game.Discard.PopFront(len(game.Discard.List) - 1)
		game.Deck.Push(recyclable...)
		game.Deck.Shuffle()
	}

	if count > len(game.Deck.List) {
		// TODO: Handle
		panic("Too many players. Not enough cards.")
	}

	cards := game.Deck.PopN(count)
	hand.Push(cards...)
}
