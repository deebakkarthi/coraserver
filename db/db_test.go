package db

import (
	"fmt"
	"testing"
)

func TestGetAllSlot(t *testing.T) {
	result := GetAllSlot()
	correct := []int{1, 2, 3, 4, 5, 6, 7, 8}
	passed := true
	for idx, val := range result {
		if val != correct[idx] {
			passed = false
			t.Errorf("GetAllSlot() = %v; want %v\n", val, correct[idx])
		}
	}
	if passed {
		fmt.Println("PASS")
	} else {
		fmt.Println("FAIL")
	}
}

func TestGetFreeClass(t *testing.T) {
	result := GetFreeClass(8, "THU")
	if result == nil {
		fmt.Println("PASS")
	} else {
		t.Errorf(`GetFreeClass(8, "THU") = %v; want nil`, result)
	}
}
