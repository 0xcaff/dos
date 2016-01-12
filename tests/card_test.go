package test

import (
	"github.com/caffinatedmonkey/dos/models"
	"testing"
)

func TestCardCover(t *testing.T) {
	var baseCard dos.Card
	var coverCard dos.Card

	// Different Color, Same Number
	baseCard = dos.Card{
		Number:    1,
		CardType:  dos.Normal,
		CardColor: dos.Red,
	}
	coverCard = dos.Card{
		Number:    1,
		CardType:  dos.Normal,
		CardColor: dos.Green,
	}
	if baseCard.CanCover(coverCard) == false {
		t.Fail()
	}

	// Same Color, Different Number
	baseCard = dos.Card{
		Number:    1,
		CardType:  dos.Normal,
		CardColor: dos.Red,
	}
	coverCard = dos.Card{
		Number:    5,
		CardType:  dos.Normal,
		CardColor: dos.Red,
	}
	if baseCard.CanCover(coverCard) == false {
		t.Fail()
	}

	// Same Color, Same Number
	baseCard = dos.Card{
		Number:    7,
		CardType:  dos.Normal,
		CardColor: dos.Green,
	}
	coverCard = dos.Card{
		Number:    7,
		CardType:  dos.Normal,
		CardColor: dos.Green,
	}
	if baseCard.CanCover(coverCard) == false {
		t.Fail()
	}

	// Different Color, Different Number
	baseCard = dos.Card{
		Number:    3,
		CardType:  dos.Normal,
		CardColor: dos.Blue,
	}
	coverCard = dos.Card{
		Number:    7,
		CardType:  dos.Normal,
		CardColor: dos.Red,
	}
	if baseCard.CanCover(coverCard) == true {
		t.Fail()
	}

	// TODO: More Scenarios
}
