package main

import "container/heap"

type CostState struct {
	mode EncodingMode
	cost int
}
type PriorityQueue []CostState

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(CostState))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
func precomputeMinimalCosts() (map[EncodingMode]map[EncodingMode]int, map[EncodingMode]map[EncodingMode]int) {
	modes := []EncodingMode{Upper, Lower, Mixed, Punct, Digit, Binary}
	shiftCost := make(map[EncodingMode]map[EncodingMode]int)
	latchCost := make(map[EncodingMode]map[EncodingMode]int)

	for _, src := range modes {
		shiftCost[src] = make(map[EncodingMode]int)
		latchCost[src] = make(map[EncodingMode]int)
		for _, dst := range modes {
			shiftCost[src][dst] = E
			latchCost[src][dst] = E
		}
	}

	for _, src := range modes {
		distShift := make(map[EncodingMode]int)
		distLatch := make(map[EncodingMode]int)
		for _, mode := range modes {
			distShift[mode] = E
			distLatch[mode] = E
		}
		distShift[src] = 0
		distLatch[src] = 0

		pq := &PriorityQueue{}
		heap.Init(pq)
		heap.Push(pq, CostState{src, 0})

		for pq.Len() > 0 {
			current := heap.Pop(pq).(CostState)
			currentMode := current.mode
			currentCost := current.cost

			for nextMode, costs := range changeLen[currentMode] {
				if costs.Shift != E {
					newCost := currentCost + costs.Shift
					if newCost < distShift[nextMode] {
						distShift[nextMode] = newCost
						heap.Push(pq, CostState{nextMode, newCost})
					}
				}
				if costs.Latch != E {
					newCost := currentCost + costs.Latch
					if newCost < distLatch[nextMode] {
						distLatch[nextMode] = newCost
						heap.Push(pq, CostState{nextMode, newCost})
					}
				}
			}
		}

		for _, dst := range modes {
			shiftCost[src][dst] = distShift[dst]
			latchCost[src][dst] = distLatch[dst]
		}
	}

	return shiftCost, latchCost
}
