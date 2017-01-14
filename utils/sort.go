package utils

import "sort"

type SortLess func(i, j int) bool
type SortSwap func(i, j int)
type sortInterfaceProxy struct {
	length int
	less   SortLess
	swap   SortSwap
}

func (p sortInterfaceProxy) Len() int           { return p.length }
func (p sortInterfaceProxy) Less(i, j int) bool { return p.less(i, j) }
func (p sortInterfaceProxy) Swap(i, j int)      { p.swap(i, j) }

func SortInterface(length int, less SortLess, swap SortSwap) sort.Interface {
	return &sortInterfaceProxy{
		length: length,
		less:   less,
		swap:   swap,
	}
}
