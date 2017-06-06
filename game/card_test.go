package dos

import (
	proto "github.com/0xcaff/dos/proto"
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

func TestDifferentSpecialCover(t *testing.T) {
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

func TestSpecialCover(t *testing.T) {
	baseCard := &proto.Card{
		Type:  proto.CardType_DOUBLEDRAW,
		Color: proto.CardColor_RED,
	}

	coverCard := &proto.Card{
		Type:  proto.CardType_DOUBLEDRAW,
		Color: proto.CardColor_GREEN,
	}

	if CanCoverCard(baseCard, coverCard) != true {
		t.Fail()
	}
}

func TestBlackCovering(t *testing.T) {
	baseCard := &proto.Card{
		Type:  proto.CardType_WILD,
		Color: proto.CardColor_BLACK,
	}

	coverCard := &proto.Card{
		Type:  proto.CardType_DOUBLEDRAW,
		Color: proto.CardColor_GREEN,
	}

	if CanCoverCard(baseCard, coverCard) != true {
		t.Fail()
	}
}

func TestBlackCover(t *testing.T) {
	baseCard := &proto.Card{
		Type:  proto.CardType_QUADDRAW,
		Color: proto.CardColor_BLACK,
	}

	coverCard := &proto.Card{
		Type:  proto.CardType_DOUBLEDRAW,
		Color: proto.CardColor_GREEN,
	}

	if CanCoverCard(coverCard, baseCard) != true {
		t.Fail()
	}
}

func TestCoveringSetWild(t *testing.T) {
	baseCard := &proto.Card{
		Type:  proto.CardType_WILD,
		Color: proto.CardColor_GREEN,
	}

	coverCard := &proto.Card{
		Type:   proto.CardType_NORMAL,
		Color:  proto.CardColor_GREEN,
		Number: 5,
	}

	if CanCoverCard(coverCard, baseCard) != true {
		t.Fail()
	}
}
