package dos

import (
	"errors"
	"math/rand"
	"sync"

	proto "github.com/caffinatedmonkey/dos/proto"
)

type Player struct {
	Name string
	Cards
	TurnDone chan struct{}
}

type Game struct {
	playerMutex sync.Mutex
	players     []*Player

	Discard Cards
	Deck    Cards

	PlayerJoined chan string
	PlayerLeft   chan string

	currentPlayerIndex int
	isReversed         bool

	// Whether the action card of the action card at the top of the discard pile
	// has been executed.
	lastCardPlayed bool
}

// Creates a new game an initalizes its values.
func NewGame(withChannels bool) *Game {
	// Initalize Values
	g := Game{
		players:            []*Player{},
		Deck:               *PlayingDeck(),
		currentPlayerIndex: -1,
		// lastPlayerPlayed:   true,
	}

	if withChannels {
		g.PlayerJoined = make(chan string)
		g.PlayerLeft = make(chan string)
	}

	g.Discard.Push(g.Deck.Pop())
	return &g
}

// Creates a new player, populates its deck and adds it to the game, notifying
// PlayerJoined.
func (game *Game) NewPlayer(name string) (*Player, error) {
	for _, player := range game.players {
		if player.Name == name {
			return nil, errors.New("Player already exists in game")
		}
	}

	player := &Player{
		Cards:    *NewCardCollection(),
		Name:     name,
		TurnDone: make(chan struct{}, 1),
	}
	game.DrawCards(&player.Cards, 8)

	game.playerMutex.Lock()
	game.players = append(game.players, player)
	game.playerMutex.Unlock()

	// Inform Listeners
	if game.PlayerJoined != nil {
		game.PlayerJoined <- name
	}

	return player, nil
}

// If player in game, removes them from the game and send a message to the
// PlayerLeft channel.
func (game *Game) RemovePlayer(removing *Player) {
	game.playerMutex.Lock()

	// Remove Player
	i := 0
	newPlayers := make([]*Player, len(game.players))
	for _, player := range game.players {
		if player != removing {
			newPlayers[i] = player
			i++
		}
	}

	newPlayers = newPlayers[:i]

	// Slice off extras
	game.players = newPlayers
	game.playerMutex.Unlock()

	// Return to discard pile
	game.Discard.PushFront(removing.Cards.List...)

	// Notify Players
	if game.PlayerLeft != nil {
		game.PlayerLeft <- removing.Name
	}
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
// next. Returns nil if there aren't enough players.
func (game *Game) NextPlayer() *Player {
	if len(game.players) < 2 {
		return nil
	}

	// If this is the first turn, pick a random player
	if game.currentPlayerIndex == -1 {
		game.currentPlayerIndex = rand.Intn(len(game.players) - 1)
		return game.players[game.currentPlayerIndex]
	}

	increment := 1

	if !game.lastCardPlayed {
		// Don't handle special actions more than once.
		lastCard := game.Discard.List[len(game.Discard.List)-1]

		switch lastCard.Type {
		case proto.CardType_REVERSE:
			game.isReversed = !game.isReversed
			if len(game.players) == 2 {
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

		game.lastCardPlayed = true
	}

	nextPlayer, index := game.GetPlayer(increment)
	game.currentPlayerIndex = index
	return nextPlayer
}

// Gets the player index n positions away from the current player.
func (game *Game) GetPlayer(n int) (*Player, int) {
	if game.isReversed {
		n = -n
	}

	playersCount := len(game.players)

	// This modulo operator isn't the same as the mathematical modulo operator.
	// It returns negative values when the left side is negative.
	current := (game.currentPlayerIndex + n) % playersCount
	if current < 0 {
		current += playersCount
	}

	return game.players[current], current
}

func (game *Game) PlayCard(player *Player, id int32, color proto.CardColor) error {
	card, index := player.Cards.FindById(id)
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

	game.lastCardPlayed = false
	game.Discard.Push(*card)
	player.Cards.RemoveCard(index)
	return nil
}

// Draws count cards into hand. Recycle's cards if needed
func (game *Game) DrawCards(hand *Cards, count int) {
	if count > len(game.Deck.List) {
		// Not enough cards to draw. Recycle entire discard pile.
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
