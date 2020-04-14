package util

import (
    "fmt"
    "testing"
)

func TestIntEqual(t *testing.T) {
    slice1 := []int{1, 2, 4, 5}
    slice2 := []int{8, 9, 4, 5}
    fmt.Print(IntEqual(slice1, slice2))
}
