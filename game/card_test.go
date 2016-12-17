package dos

import (
	proto "github.com/caffinatedmonkey/dos/proto"
	"testing"
)

func TestCoverSameNumber(t *testing.T) {
	baseCard := &proto.Card{
		Number: 1,
		Type:   proto.CardType_NORMAL,
		Color:  proto.CardColor_RED,
	}
	coverCard := &proto.Card{
		Number: 1,
		Type:   proto.CardType_NORMAL,
		Color:  proto.CardColor_GREEN,
	}

	if CanCoverCard(baseCard, coverCard) != true {
		t.Fail()
	}
}

func TestCoverSameColor(t *testing.T) {
	baseCard := &proto.Card{
		Number: 1,
		Type:   proto.CardType_NORMAL,
		Color:  proto.CardColor_RED,
	}
	coverCard := &proto.Card{
		Number: 5,
		Type:   proto.CardType_NORMAL,
		Color:  proto.CardColor_RED,
	}

	if CanCoverCard(baseCard, coverCard) != true {
		t.Fail()
	}
}

func TestCoverSameColorNumber(t *testing.T) {
	// Same Color, Same Number
	baseCard := &proto.Card{
		Number: 7,
		Type:   proto.CardType_NORMAL,
		Color:  proto.CardColor_GREEN,
	}

	coverCard := &proto.Card{
		Number: 7,
		Type:   proto.CardType_NORMAL,
		Color:  proto.CardColor_GREEN,
	}

	if CanCoverCard(baseCard, coverCard) != true {
		t.Fail()
	}
}

func TestCoverDifferentColor(t *testing.T) {
	baseCard := &proto.Card{
		Number: 3,
		Type:   proto.CardType_NORMAL,
		Color:  proto.CardColor_BLUE,
	}

	coverCard := &proto.Card{
		Number: 7,
		Type:   proto.CardType_NORMAL,
		Color:  proto.CardColor_RED,
	}

	if CanCoverCard(baseCard, coverCard) != false {
		t.Fail()
	}

}

func TestSpecialCover(t *testing.T) {
	baseCard := &proto.Card{
		Type:  proto.CardType_REVERSE,
		Color: proto.CardColor_RED,
	}

	coverCard := &proto.Card{
		Type:  proto.CardType_DOUBLEDRAW,
		Color: proto.CardColor_GREEN,
	}

	if CanCoverCard(baseCard, coverCard) != false {
		t.Fail()
	}
}
