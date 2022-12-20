package websockets

import (
	"reflect"
	"testing"
)

func Test_removeFromSlice(t *testing.T) {
	// Test removing a value from the middle of the slice
	slice := []int{1, 2, 3, 4, 5}
	slice = removeFromSlice(slice, 2)
	expected := []int{1, 2, 4, 5}
	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("RemoveFromSlice: expected %v, got %v", expected, slice)
	}

	// Test removing the first value in the slice
	slice = []int{1, 2, 3, 4, 5}
	slice = removeFromSlice(slice, 0)
	expected = []int{2, 3, 4, 5}
	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("RemoveFromSlice: expected %v, got %v", expected, slice)
	}

	// Test removing the last value in the slice
	slice = []int{1, 2, 3, 4, 5}
	slice = removeFromSlice(slice, 4)
	expected = []int{1, 2, 3, 4}
	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("RemoveFromSlice: expected %v, got %v", expected, slice)
	}

	// Test removing a value from an index out of bounds
	slice = []int{1, 2, 3, 4, 5}
	slice = removeFromSlice(slice, 5)
	expected = []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("RemoveFromSlice: expected %v, got %v", expected, slice)
	}
}
