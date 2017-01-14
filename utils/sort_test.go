package utils

import (
	"sort"
	"testing"
)

var arr = []int{1, 2, 3}
var si = SortInterface(len(arr), func(i, j int) bool {
	return i < j
}, func(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
})

func TestReverse(t *testing.T) {
	sort.Sort(sort.Reverse(si))
	//t.Fatal(arr)
}

func BenchmarkReverse(b *testing.B) {
	sort.Sort(sort.Reverse(si))
}
