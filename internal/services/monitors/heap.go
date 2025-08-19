package monitors

import (
	"container/heap"
	"time"
)

type taskHeap []*ScheduledTask

func (h taskHeap) Len() int { return len(h) }
func (h taskHeap) Less(i, j int) bool {
	return h[i].NextCheckAt.Before(h[j].NextCheckAt)
}
func (h taskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *taskHeap) Pop() any {
	x := (*h)[len(*h)-1]
	x.index = -1
	*h = (*h)[:len(*h)-1]
	return x
}

func (h *taskHeap) Push(x any) {
	task := x.(*ScheduledTask)
	task.index = len(*h)
	*h = append(*h, task)
}

func (h *taskHeap) update(task *ScheduledTask, next time.Time) {
	task.NextCheckAt = next
	heap.Fix(h, task.index)
}
