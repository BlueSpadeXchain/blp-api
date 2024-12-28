package taskqueue

import (
	"container/heap"
)

type Task struct {
	Timestamp  int64
	IsBalancer bool
	Data       interface{} // This can be the data you need to process
}

// we need to figure out processing one order
// then we need to modify this to execute one or many orders
type PriorityQueue []*Task

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].Timestamp == pq[j].Timestamp {
		return pq[i].IsBalancer && !pq[j].IsBalancer
	}
	return pq[i].Timestamp < pq[j].Timestamp
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Task))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func ProcessQueue(pq *PriorityQueue) {
	for pq.Len() > 0 {
		task := heap.Pop(pq).(*Task)
		// Process the task here
		if task.IsBalancer {
			processBalancerTask(task)
		} else {
			processOrderTask(task)
		}
	}
}

func processBalancerTask(task *Task) {
	// Your logic to handle balancer bot tasks
}

func processOrderTask(task *Task) {
	// Your logic to handle order tasks
}
