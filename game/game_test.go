package dos

import (
	proto "github.com/0xcaff/dos/proto"
	"testing"
)

func SetupGoldilocksGame() (*Game, *Player, *Player, *Player, *Player) {
	game := NewGame(false)

	littleBear, _ := game.NewPlayer("Little, Small, Wee Bear")
	mediumBear, _ := game.NewPlayer("Medium-sized Bear")
	hugeBear, _ := game.NewPlayer("Great, Huge Bear")
	goldilocks, _ := game.NewPlayer("Goldilocks")

	game.currentPlayerIndex = 0

	return game, littleBear, mediumBear, hugeBear, goldilocks
}

func TestNextPlayer(t *testing.T) {
	game, littleBear, mediumBear, hugeBear, goldilocks := SetupGoldilocksGame()

	for i := 0; i < 20; i++ {
		nextPlayer := game.NextPlayer()
		if nextPlayer != mediumBear {
			t.Fail()
		}

		nextPlayer = game.NextPlayer()
		if nextPlayer != hugeBear {
			t.Fail()
		}

		nextPlayer = game.NextPlayer()
		if nextPlayer != goldilocks {
			t.Fail()
		}

		nextPlayer = game.NextPlayer()
		if nextPlayer != littleBear {
			t.Fail()
		}
	}
}

func TestGetPlayer(t *testing.T) {
	game, _, _, _, _ := SetupGoldilocksGame()

	nextValue := 0
	for i := -100; i < 100; i++ {
		_, index := game.GetPlayer(i)
		if index < 0 || index > 3 || nextValue != index {
			t.Fail()
		}
		nextValue++

		if nextValue == 4 {
			nextValue = 0
		}
	}
}

func TestExpendCards(t *testing.T) {
	game := NewGame(false)

	lonePlayer, _ := game.NewPlayer("A lonely player")
	game.DrawCards(&game.Discard, 99)
	game.DrawCards(&lonePlayer.Cards, 20)
}

func TestPlayCards(t *testing.T) {
	// Setup game
	game := NewGame(false)
	game.Discard.Push(proto.Card{
		Color:  proto.CardColor_RED,
		Type:   proto.CardType_NORMAL,
		Number: int32(5),
	})

	toPlay := proto.Card{
		Color:  proto.CardColor_RED,
		Type:   proto.CardType_NORMAL,
		Number: int32(10),
		Id:     int32(200), // Nothing will normally ever get an id this high.
	}

	lonePlayer, _ := game.NewPlayer("A lonely player")

	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			lonePlayer.Cards.PushFront(toPlay)
		} else {
			lonePlayer.Cards.Push(toPlay)
		}

		// Try playing it
		err := game.PlayCard(lonePlayer, toPlay.Id, proto.CardColor(0))
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		// Check whether it is played
		discard := game.Discard.List[len(game.Discard.List)-1]
		if discard.Id != toPlay.Id {
			t.Fail()
		}

		toPlay.Id++
	}
}

func TestDuplicateSpecial(t *testing.T) {
	game := NewGame(false)
	game.currentPlayerIndex = -1

	player1, _ := game.NewPlayer("Player 1")
	player2, _ := game.NewPlayer("Player 2")

	if player := game.NextPlayer(); player1 != player {
		t.Log("First Turn", player1, player)
		t.Fail()
	}

	game.Discard.Push(proto.Card{
		Color: proto.CardColor_RED,
		Type:  proto.CardType_SKIP,
	})

	if player := game.NextPlayer(); player1 != player {
		t.Log("Second Turn. Expected:", player1, "Got:", player)
		t.Fail()
	}

	if player := game.NextPlayer(); player2 != player {
		t.Log("Third Turn. Expected:", player2, "Got:", player)
		t.Fail()
	}

	if player := game.NextPlayer(); player1 != player {
		t.Log("Fourth Turn. Expected:", player1, "Got:", player)
		t.Fail()
	}
}

func TestRemovePlayer(t *testing.T) {
	game, littleBear, mediumBear, hugeBear, goldilocks := SetupGoldilocksGame()

	game.RemovePlayer(littleBear)
	game.RemovePlayer(goldilocks)

	if len(game.players) != 2 {
		t.Log("Invalid number of players. Expected 2. Got: ", len(game.players))
		t.Fail()
	}

	for i := 0; i < 20; i++ {
		player := game.NextPlayer()
		if player != mediumBear && player != hugeBear {
			t.Log("Invalid player", player)
			t.Fail()
		}
	}
}

func TestNextPlayerEmpty(t *testing.T) {
	game := NewGame(false)
	player := game.NextPlayer()
	if player != nil {
		t.Log("Player in empty game", player)
		t.Fail()
	}
}

func TestRemoveOtherPlayer(t *testing.T) {
	game, _, _, _, _ := SetupGoldilocksGame()

	anotherRandomPlayer := &Player{Name: "Bruce"}
	game.RemovePlayer(anotherRandomPlayer)

	players := len(game.players)
	if players != 4 {
		t.Log("Invalid number of players. Expected 4. Got: ", players)
		t.Fail()
	}
}

// TODO: Test Special Card Actions
