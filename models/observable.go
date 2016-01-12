package dos

import "math/rand"

type ObservableList struct {
	List []interface{}

	additions chan interface{}
	deletions chan interface{}
}

func (ol *ObservableList) Push(e ...interface{}) {
	for _, elem := range e {
		if ol.additions != nil {
			ol.additions <- elem
		}
		ol.List = append(ol.List, elem)
	}
}

// Removes and returns the last element of the array. If the array is empty,
// returns nil.
func (ol *ObservableList) Pop() interface{} {
	lastElementIndex := len(ol.List) - 1
	elem := ol.List[lastElementIndex]
	ol.List = ol.List[:lastElementIndex]
	if ol.deletions != nil {
		ol.deletions <- elem
	}
	return elem
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

// TODO: Make this respond differently in case of failure
func (ol *ObservableList) GetElement(index int) interface{} {
	if index <= len(ol.List)-1 {
		return ol.List[index]
	} else {
		return nil
	}
}

// Calls f with element when an element is added.
func (ol *ObservableList) OnAddition(f func(interface{})) {
	registerListener(f, &ol.additions)
}

// Calls f with element when an element is deleted.
func (ol *ObservableList) OnDeletion(f func(interface{})) {
	registerListener(f, &ol.deletions)
}

func (ol *ObservableList) Length() int {
	return len(ol.List)
}

// Calls function f when a value is pushed onto channel c.
func registerListener(f func(interface{}), c *chan interface{}) {
	if *c == nil {
		*c = make(chan interface{})
	}
	go func() {
		for e := range *c {
			f(e)
		}
	}()
}

// Swaps elements at array[f] and array[t] with each other.
func swap(f, t int, array []interface{}) {
	fl := array[f]
	tl := array[t]
	array[t] = fl
	array[f] = tl
}
