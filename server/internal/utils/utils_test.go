package utils

import "testing"

func TestRangeNil(t *testing.T) {
	var nilMap map[string]int = nil
	for k, v := range nilMap {
		t.Log(k, v)
	}
}

func TestNumberSlice(t *testing.T) {
	slice := []int{1}

	t.Log(slice[0], slice[1:])
}
