package main

import (
	"sort"
	"sync"
)

type queue struct {
	sync.Mutex
	tasks []*task
}

func newQueue() *queue {
	return &queue{
		tasks: []*task{},
	}
}

func (q *queue) next() (t *task) {
	q.Lock()
	defer q.Unlock()
	if len(q.tasks) > 0 {
		t = q.tasks[0]
		q.tasks = q.tasks[1:]
	}
	return t
}

func (q *queue) add(t *task) {
	q.Lock()
	defer q.Unlock()

	q.tasks = append(q.tasks, t)

	// sort for 1. prio asc and 2. start time asc
	sort.Slice(q.tasks, func(i, j int) bool {
		if q.tasks[i].Prio < q.tasks[j].Prio {
			return true
		}
		if q.tasks[i].Prio > q.tasks[j].Prio {
			return false
		}
		return q.tasks[i].Start.Before(q.tasks[j].Start)
	})
}
