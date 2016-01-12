package test

import (
	"github.com/caffinatedmonkey/dos/models"
	"testing"
)

// TODO: This test will fail if even one of the elements is not shuffled.
func TestShuffle(t *testing.T) {
	l := dos.ObservableList{}
	l.Push(1, 2, 3, 4, 5)
	l.Shuffle()
	for i := 0; i < 5; i++ {
		elem := l.GetElement(i)
		if elem == nil || elem.(int) == i {
			t.Logf("Failed (shuffled) %d == (nonshuffled) %d", elem, i)
			t.Fail()
		} else {
			t.Logf("Passed (shuffled) %d == (nonshuffled) %d", elem, i)
		}
	}
}

func TestPushing(t *testing.T) {
	l := dos.ObservableList{}
	l.Push(1)
	initialLength := l.Length()

	l.Push(1, 2, 3)
	finalLength := l.Length()

	if initialLength == finalLength {
		t.Fail()
	}
}

// TODO: No failing situtation
func TestObservation(t *testing.T) {
	l := dos.ObservableList{}
	l.OnAddition(func(i interface{}) {
		t.Logf("Added %d", i)
	})
	l.Push(1, 2, 3)
}
