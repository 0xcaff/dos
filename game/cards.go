package dos

import (
	"math/rand"
	"sync"

	proto "github.com/caffinatedmonkey/dos/proto"
	"github.com/caffinatedmonkey/dos/utils"
)

type Cards struct {
	List []proto.Card

	Deletions utils.Broadcaster
	Additions utils.Broadcaster

	sync.RWMutex
}

func NewCardCollection() *Cards {
	return &Cards{
		List:      []proto.Card{},
		Additions: *utils.NewBroadcaster(),
		Deletions: *utils.NewBroadcaster(),
	}
}

func (ol *Cards) Push(e ...proto.Card) {
	ol.RWMutex.Lock()
	defer ol.RWMutex.Unlock()

	ol.List = append(ol.List, e...)
	for _, elem := range e {
		ol.Additions.Broadcast(elem)
	}
}

func (cards *Cards) PushFront(newCards ...proto.Card) {
	cards.RWMutex.Lock()
	defer cards.RWMutex.Unlock()

	cards.List = append(newCards, cards.List...)
	for _, card := range newCards {
		cards.Additions.Broadcast(card)
	}
}

// Removes and returns the last element of the array. If the array is empty,
// returns nil.
func (ol *Cards) Pop() proto.Card {
	return ol.PopN(1)[0]
}

// Removes and returns the last n elements of the array.
func (ol *Cards) PopN(n int) []proto.Card {
	ol.RWMutex.Lock()
	defer ol.RWMutex.Unlock()

	lengthBefore := len(ol.List)
	lastElementIndex := lengthBefore - n
	if lastElementIndex < 0 {
		return nil
	}

	// Get Last n Elements
	elems := ol.List[lastElementIndex:]

	for index := range elems {
		ol.Deletions.Broadcast(elems[index].Id)
	}

	// Remove Last n Elements
	ol.List = ol.List[:lastElementIndex]

	return elems
}

func (cards *Cards) PopFront(n int) []proto.Card {
	if n > len(cards.List) {
		return nil
	}

	cards.RWMutex.Lock()
	defer cards.RWMutex.Unlock()

	removing := cards.List[:n]
	cards.List = cards.List[n:]

	for index := range removing {
		cards.Deletions.Broadcast(removing[index].Id)
	}

	return removing
}

// Shuffles the underlying list using a fisher yates shuffle.
func (cards *Cards) Shuffle() {
	cards.RWMutex.Lock()
	defer cards.RWMutex.Unlock()

	currentIndex := len(cards.List)

	for currentIndex > 0 {
		randomIndex := rand.Intn(currentIndex)
		currentIndex = currentIndex - 1
		swap(currentIndex, randomIndex, cards.List)
	}
}

func (cards *Cards) FindById(id int32) (*proto.Card, int) {
	cards.RWMutex.RLock()
	defer cards.RWMutex.RUnlock()
	for index := range cards.List {
		card := &cards.List[index]
		if card.Id == id {
			return card, index
		}
	}

	return nil, -1
}

func (cards *Cards) RemoveCard(index int) {
	cards.RWMutex.Lock()
	defer cards.RWMutex.Unlock()

	deleting := cards.List[index]
	cards.List = append(cards.List[:index], cards.List[index+1:]...)
	cards.Deletions.Broadcast(deleting.Id)
}

// func (cards *Cards) PopId(id int32) *proto.Card {
// 	cards.RWMutex.Lock()
// 	defer cards.RWMutex.Unlock()
//
// 	newCards := make([]proto.Card, len(cards.List))
// 	var foundCard *proto.Card
// 	i := 0
// 	for index := range cards.List {
// 		card := cards.List[index]
// 		if card.Id == id {
// 			foundCard = &card
// 		} else {
// 			newCards[i] = card
// 			i++
// 		}
// 	}
//
// 	cards.List = newCards
// 	if foundCard != nil {
// 		cards.Deletions.Broadcast(foundCard.Id)
// 	}
//
// 	return foundCard
// }

func swap(f, t int, array []proto.Card) {
	array[t], array[f] = array[f], array[t]
}

// func (ol *Cards) GetElement(index int) interface{} {
// 	if index <= len(ol.list)-1 {
// 		return ol.list[index]
// 	} else {
// 		panic("Index out of bounds")
// 	}
// }

// Creates a new deck and populates it with the standard playing cards
func PlayingDeck() *Cards {
	deck := NewCardCollection()
	deck.Populate()
	deck.Shuffle()
	return deck
}

var CardColors = []proto.CardColor{
	proto.CardColor_RED,
	proto.CardColor_YELLOW,
	proto.CardColor_GREEN,
	proto.CardColor_BLUE,
}

// Add the standard cards to the deck.
func (cards *Cards) Populate() {
	id := int32(0)
	for i := 0; i < 2; i++ {
		for _, color := range CardColors {
			// Insert Cards 1-9
			for k := int32(1); k < int32(10); k++ {
				id++
				cards.Push(proto.Card{
					Id:     id,
					Number: k,
					Color:  color,
					Type:   proto.CardType_NORMAL,
				})
			}

			// Insert Special Cards (Skip, Reverse, DoubleDraw)
			id++
			cards.Push(proto.Card{
				Id:     id,
				Number: -1,
				Color:  color,
				Type:   proto.CardType_SKIP,
			})

			id++
			cards.Push(proto.Card{
				Id:     id,
				Number: -1,
				Color:  color,
				Type:   proto.CardType_REVERSE,
			})

			id++
			cards.Push(proto.Card{
				Id:     id,
				Number: -1,
				Color:  color,
				Type:   proto.CardType_DOUBLEDRAW,
			})

			if i == 0 {
				id++
				cards.Push(proto.Card{
					Id:     id,
					Number: -1,
					Color:  proto.CardColor_BLACK,
					Type:   proto.CardType_QUADDRAW,
				})

				id++
				cards.Push(proto.Card{
					Id:     id,
					Number: 0,
					Color:  color,
					Type:   proto.CardType_NORMAL,
				})
			} else if i == 1 {
				id++
				cards.Push(proto.Card{
					Id:     id,
					Number: -1,
					Color:  proto.CardColor_BLACK,
					Type:   proto.CardType_WILD,
				})
			}
		}
	}
}

func CanCoverCard(baseCard, otherCard *proto.Card) bool {
	specialsMatch := baseCard.Type == otherCard.Type
	bothAreNormal := specialsMatch && baseCard.Type == proto.CardType_NORMAL

	colorsMatch := baseCard.Color == otherCard.Color
	numbersMatch := baseCard.Number == otherCard.Number

	if colorsMatch || (numbersMatch && bothAreNormal) {
		return true
	}

	bothAreNotNormal := specialsMatch && baseCard.Type != proto.CardType_NORMAL
	if specialsMatch && bothAreNotNormal {
		return true
	}

	if baseCard.Color == proto.CardColor_BLACK {
		// If the game starts with a black card, anything can cover it.
		return true
	}

	return false
}
