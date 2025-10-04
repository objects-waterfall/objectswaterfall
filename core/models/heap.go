package models

import (
	"container/heap"
)

type MedianValue struct {
	low  *MaxHeap
	high *MinHeap
}

type MaxHeap []float64
type MinHeap []float64

func NewMedianValue() MedianValue {
	l := &MaxHeap{}
	h := &MinHeap{}
	heap.Init(l)
	heap.Init(h)
	return MedianValue{low: l, high: h}
}

func (mf *MedianValue) AddNum(num float64) {
	heap.Push(mf.low, num)
	heap.Push(mf.high, heap.Pop(mf.low))

	if mf.low.Len() < mf.high.Len() {
		heap.Push(mf.low, heap.Pop(mf.high))
	}
}

func (mf *MedianValue) FindMedian() float64 {
	if mf.low.Len() > mf.high.Len() {
		return float64((*mf.low)[0])
	}
	return float64((*mf.low)[0]+(*mf.high)[0]) / 2.0
}

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(float64)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(float64)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
