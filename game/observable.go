package dos

import (
	"math/rand"
	"sync"
)

type ObservableList struct {
	List []interface{}

	// Objects are accumulated in these buffers as events trigger. When a
	// listener is attached all objects stored in the buffer call the listener.
	// Once flush is called, these bufferes are emptied.
	additionBuffer []interface{}
	deletionBuffer []int

	// These channels are triggered when one or more objects is added or deleted
	// from one of the buffers. These channels are not triggered once flushed.
	onAddition chan struct{}
	onDeletion chan struct{}

	additionCallbacks []func(...interface{})
	deletionCallbacks []func(...int)

	sync.Mutex
}

func NewObservableList() ObservableList {
	return ObservableList{
		List: []interface{}{},
	}
}

func (ol *ObservableList) Push(e ...interface{}) {
	ol.Mutex.Lock()

	ol.List = append(ol.List, e...)

	ol.additionBuffer = append(ol.additionBuffer, e...)
	if ol.onAddition != nil {
		ol.onAddition <- struct{}{}
	}

	ol.Mutex.Unlock()
}

// Removes and returns the last element of the array. If the array is empty,
// returns nil.
func (ol *ObservableList) Pop() interface{} {
	return ol.PopN(1)[0]
}

// Removes and returns the last n elements of the array.
func (ol *ObservableList) PopN(int n) []interface{} {
	lengthBefore := len(ol.List)
	lastElementIndex := lengthBefore - n
	if lastElementIndex < 0 {
		panic("Popping Too Many Elements")
	}

	ol.Mutex.Lock()

	// Get Last n Elements
	elems := ol.List[lastElementIndex:]

	// Remove Last n Elements
	ol.List = ol.List[:lastElementIndex]

	// Notify Listeners
	var deletions [lenBefore - len(ol.List)]int
	for i := lastElementIndex; i < lenBefore; i++ {
		deletions[i-lastElementIndex] = i
	}

	ol.deletionBuffer = append(ol.deletionBuffer, deletions)
	if ol.onDeletion != nil {
		ol.onDeletion <- struct{}{}
	}

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
	re := ol.GetElement(index)
	if ol.deletions != nil {
		ol.deletions <- re
	}
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

// Calls f with element when an element is added.
func (ol *ObservableList) OnAddition(f func(interface{})) {
	if len(ol.additionCallbacks) == 0 {
		// First Callback, Register Dispatcher

		go func(ol *ObservableList) {
			for event := range ol.onAddition {

				// Dispatch Events
				for cb := range ol.additionCallbacks {
					cb(ol.additionBuffer...)
				}

				ol.additionBuffer = nil
			}
		}(ol)
	}

	// Register Listener
	ol.additionCallbacks = append(ol.additionCallbacks, f)
	f(ol.additionBuffer...)

	// TODO: Race Condition
}

// Calls f with element when an element is deleted.
func (ol *ObservableList) OnDeletion(f func(int)) {
	// TODO: Generalize
}

// Swaps elements at array[f] and array[t] with each other.
func swap(f, t int, array []interface{}) {
	fl := array[f]
	tl := array[t]
	array[t] = fl
	array[f] = tl
}
