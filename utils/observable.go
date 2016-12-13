package utils

import (
	"math/rand"
	"sync"
)

type ObservableList struct {
	List []interface{}

	Deletions Broadcaster
	Additions Broadcaster

	sync.Mutex
}

func NewObservableList() ObservableList {
	return ObservableList{
		List:      []interface{}{},
		Additions: *NewBroadcaster(),
		Deletions: *NewBroadcaster(),
	}
}

func (ol *ObservableList) Push(e ...interface{}) {
	ol.Mutex.Lock()
	ol.List = append(ol.List, e...)

	for _, elem := range e {
		ol.Additions.Broadcast(elem)
	}
	ol.Mutex.Unlock()
}

// Removes and returns the last element of the array. If the array is empty,
// returns nil.
func (ol *ObservableList) Pop() interface{} {
	return ol.PopN(1)[0]
}

// Removes and returns the last n elements of the array.
func (ol *ObservableList) PopN(n int) []interface{} {
	lengthBefore := len(ol.List)
	lastElementIndex := lengthBefore - n
	if lastElementIndex < 0 {
		panic("Popping Too Many Elements")
	}

	ol.Mutex.Lock()

	// Get Last n Elements
	elems := ol.List[lastElementIndex:]

	for i := lastElementIndex; i <= lastElementIndex+n; i++ {
		ol.Deletions.Broadcast(i)
	}

	// Remove Last n Elements
	ol.List = ol.List[:lastElementIndex]

	ol.Mutex.Unlock()
	return elems
}

// Shuffles the underlying list using a fisher yates shuffle.
func (ol *ObservableList) Shuffle() {
	currentIndex := len(ol.List)

	for currentIndex > 0 {
		randomIndex := rand.Intn(currentIndex)
		currentIndex = currentIndex - 1
		swap(currentIndex, randomIndex, ol.List)
	}
}

func (ol *ObservableList) GetElement(index int) interface{} {
	if index <= len(ol.List)-1 {
		return ol.List[index]
	} else {
		panic("Index out of bounds")
	}
}

func (ol *ObservableList) RemoveElement(index int) {
	ol.Mutex.Lock()
	ol.Deletions.Broadcast(index)
	ol.List = append(ol.List[:index], ol.List[index+1:]...)
	ol.Mutex.Unlock()
}

func (ol *ObservableList) Clear() {
	ol.List = []interface{}{}
}

func (ol *ObservableList) Length() int {
	return len(ol.List)
}

func (ol *ObservableList) Each(f func(interface{}, int) bool) bool {
	for i, v := range ol.List {
		r := f(v, i)
		if r {
			return false
		}
	}
	return true
}

// Runs f on each element returning the return value for each call.
func (ol *ObservableList) Map(f func(interface{}) interface{}) []interface{} {
	r := make([]interface{}, ol.Length())
	for i, v := range ol.List {
		r[i] = f(v)
	}
	return r
}

// Swaps at array[f] and array[t]
func swap(f, t int, array []interface{}) {
	array[t], array[f] = array[f], array[t]
}
