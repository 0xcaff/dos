package dos

// import (
// 	"github.com/caffinatedmonkey/dos/models"
// 	"testing"
// )
//
// // TODO: This test will fail if even one of the elements is not shuffled.
// func TestShuffle(t *testing.T) {
// 	l := dos.NewObservableList()
// 	l.Push(1, 2, 3, 4, 5)
// 	l.Shuffle()
// 	for i := 0; i < 5; i++ {
// 		elem := l.GetElement(i)
// 		if elem == nil || elem.(int) == i {
// 			t.Logf("Failed (shuffled) %d == (nonshuffled) %d", elem, i)
// 			t.Fail()
// 		} else {
// 			t.Logf("Passed (shuffled) %d == (nonshuffled) %d", elem, i)
// 		}
// 	}
// }
//
// func TestPushing(t *testing.T) {
// 	l := dos.NewObservableList()
// 	l.Push(1)
// 	initialLength := l.Length()
//
// 	l.Push(1, 2, 3)
// 	finalLength := l.Length()
//
// 	if initialLength == finalLength {
// 		t.Fail()
// 	}
// }
//
// // TODO: No failing situtation
// func TestObservation(t *testing.T) {
// 	l := dos.NewObservableList()
// 	l.OnAddition(func(i interface{}) {
// 		t.Logf("Added %d", i)
// 	})
// 	l.Push(1, 2, 3)
// }
//
// func TestIterating(t *testing.T) {
// 	l := dos.NewObservableList()
// 	l.Push(1, 2, 3, 4, 5, 6)
// 	i := 0
//
// 	l.Each(func(v interface{}, j int) bool {
// 		i = i + 1
// 		if v.(int) != i || v.(int) != (j+1) {
// 			t.Fail()
// 		}
// 		return false
// 	})
//
// 	if i != 6 {
// 		t.Fail()
// 	}
// }
//
// func TestEarlyIteration(t *testing.T) {
// 	l := dos.NewObservableList()
// 	l.Push(1, 2, 3)
// 	i := 0
//
// 	l.Each(func(v interface{}, j int) bool {
// 		i = i + 1
// 		if v.(int) == 2 {
// 			return true
// 		} else {
// 			return false
// 		}
// 	})
//
// 	if i != 2 {
// 		t.Fail()
// 	}
// }
//
// func TestRemoveElement(t *testing.T) {
// 	l := dos.NewObservableList()
// 	l.Push(1, 2, 3)
// 	l.RemoveElement(0)
// 	if l.Length() != 2 {
// 		t.Logf("%d (length) != 2", l.Length())
// 		t.Fail()
// 	}
//
// 	l.Each(func(i interface{}, index int) bool {
// 		if in := i.(int); in != index+2 {
// 			t.Logf("%d (value) != %d", in, index+2)
// 			t.Fail()
// 		}
// 		return false
// 	})
// }
//
// func TestPop(t *testing.T) {
// 	l := dos.NewObservableList()
// 	l.Push(1, 2, 3)
// 	if i := l.Pop().(int); i != 3 {
// 		t.Logf("%d (popped) != 3", i)
// 		t.Fail()
// 	}
//
// 	l.Each(func(i interface{}, index int) bool {
// 		if in := i.(int); in != index+1 {
// 			t.Logf("%d (value) != %d (index)", in, index+1)
// 			t.Fail()
// 		}
// 		return false
// 	})
// }
//
// func TestGetElement(t *testing.T) {
// 	l := dos.NewObservableList()
// 	l.Push(1, 2, 3, 4, 5)
// 	for i := 0; i < 5; i++ {
// 		if in := l.GetElement(i).(int); in != i+1 {
// 			t.Logf("%d (got) != %d (expected)", in, i+1)
// 			t.Fail()
// 		}
// 	}
//
// 	// TODO: Fix This Case
// 	// l.GetElement(1000)
// }
//
// func TestMap(t *testing.T) {
// 	l := dos.NewObservableList()
// 	l.Push(1, 2, 3, 4)
// 	mpd := l.Map(func(i interface{}) interface{} { return i.(int) - 1 })
// 	for idx, v := range mpd {
// 		if in := v.(int); in != idx {
// 			t.Logf("%d (got) != %d (expected)", in, idx)
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestClear(t *testing.T) {
// 	l := dos.NewObservableList()
// 	l.Push(1, 2, 3, 4, 5)
// 	if l.Length() != 5 {
// 		t.Logf("%d (length) != %d (should be)", l.Length(), 5)
// 		t.Fail()
// 	}
//
// 	l.Clear()
// 	if l.Length() != 0 {
// 		t.Logf("%d (length) != %d (should be)", l.Length(), 0)
// 		t.Fail()
// 	}
// }
