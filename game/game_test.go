package dos

import (
	"testing"
)

func SetupGoldilocksGame() (*Game, *Player, *Player, *Player, *Player) {
	game := NewGame()

	littleBear, _ := game.NewPlayer("Little, Small, Wee Bear")
	mediumBear, _ := game.NewPlayer("Middle-sized Bear")
	hugeBear, _ := game.NewPlayer("Great, Huge Bear")
	goldilocks, _ := game.NewPlayer("Goldilocks")

	game.currentPlayerIndex = 0

	return game, littleBear, mediumBear, hugeBear, goldilocks
}

func TestNextPlayer(t *testing.T) {
	game, littleBear, mediumBear, hugeBear, goldilocks := SetupGoldilocksGame()

	for i := 0; i < 2; i++ {
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
	game := NewGame()

	lonePlayer, _ := game.NewPlayer("A lonely player")
	game.DrawCards(&game.Discard, 99)
	game.DrawCards(&lonePlayer.Cards, 20)
}
