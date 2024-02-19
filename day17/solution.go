package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
)

var MIN_STRAIGHT_LINES = 4
var MAX_STRAIGHT_LINES = 10

type LavaMap struct {
	costs          []int
	bestTotalCosts map[int]int
	size           int
}

func (lavaMap *LavaMap) GetCost(x int, y int) int {
	return lavaMap.costs[y*lavaMap.size+x]
}

func (lavaMap *LavaMap) getBestTotalCostKey(x int, y int, xScale int, yScale int) int {
	scaleValue := func(scale int) int {
		if scale < 0 {
			return (scale * -1) << 1
		}
		return scale
	}
	return (x << 48) | (y << 32) | (scaleValue(xScale) << 24) | (scaleValue(yScale) << 16)
}

func (lavaMap *LavaMap) GetBestTotalCost(x int, y int, xScale int, yScale int) int {
	return lavaMap.bestTotalCosts[lavaMap.getBestTotalCostKey(x, y, xScale, yScale)]
}

func (lavaMap *LavaMap) SetBestTotalCost(x int, y int, xScale int, yScale int, cost int) {
	lavaMap.bestTotalCosts[lavaMap.getBestTotalCostKey(x, y, xScale, yScale)] = cost
}

func (lavaMap *LavaMap) IsInBounds(x int, y int) bool {
	return x >= 0 && x < lavaMap.size && y >= 0 && y < lavaMap.size
}

type Path struct {
	index      int
	cost       int
	x          int
	y          int
	vertical   bool
	horizontal bool
}

// Heap implementation copied from https://pkg.go.dev/container/heap.
type PriorityQueue []*Path

func (priorityQueue PriorityQueue) Len() int { return len(priorityQueue) }

func (priorityQueue PriorityQueue) Less(i int, j int) bool {
	return priorityQueue[j].cost > priorityQueue[i].cost
}

func (priorityQueue PriorityQueue) Swap(i int, j int) {
	priorityQueue[i], priorityQueue[j] = priorityQueue[j], priorityQueue[i]
	priorityQueue[i].index = i
	priorityQueue[j].index = j
}

func (priorityQueue *PriorityQueue) Push(x any) {
	n := len(*priorityQueue)
	path := x.(*Path)
	path.index = n
	*priorityQueue = append(*priorityQueue, path)
}

func (priorityQueue *PriorityQueue) Pop() any {
	old := *priorityQueue
	n := len(old)
	path := old[n-1]
	old[n-1] = nil
	path.index = -1
	*priorityQueue = old[0 : n-1]
	return path
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lavaMap := LavaMap{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if lavaMap.size == 0 {
			lavaMap.size = len(line)
			totalSize := lavaMap.size * lavaMap.size
			lavaMap.costs = make([]int, 0, totalSize)
			lavaMap.bestTotalCosts = make(map[int]int)
		}
		for _, char := range line {
			lavaMap.costs = append(lavaMap.costs, int(char-'0'))
		}
	}

	priorityQueue := PriorityQueue{}
	heap.Push(&priorityQueue, &Path{cost: 0, x: 0, y: 0, horizontal: true, vertical: true})

	addPaths := func(path Path, xFactor int, yFactor int) {
		cost := path.cost
		for offset := 1; offset <= MAX_STRAIGHT_LINES; offset += 1 {
			xScale := xFactor * offset
			yScale := yFactor * offset
			x := path.x + xScale
			y := path.y + yScale
			if !lavaMap.IsInBounds(x, y) {
				break
			}

			cost += lavaMap.GetCost(x, y)

			if offset < MIN_STRAIGHT_LINES {
				continue
			}

			bestTotalCost := lavaMap.GetBestTotalCost(x, y, xScale, yScale)
			if bestTotalCost == 0 || cost < bestTotalCost {
				heap.Push(&priorityQueue, &Path{cost: cost, x: x, y: y, vertical: xFactor != 0, horizontal: yFactor != 0})
				lavaMap.SetBestTotalCost(x, y, xScale, yScale, cost)
			}
		}
	}

	for {
		path := heap.Pop(&priorityQueue).(*Path)
		if path.x == lavaMap.size-1 && path.y == lavaMap.size-1 {
			fmt.Println(path.cost)
			break
		}
		if path.horizontal {
			addPaths(*path, 1, 0)
			addPaths(*path, -1, 0)
		}
		if path.vertical {
			addPaths(*path, 0, 1)
			addPaths(*path, 0, -1)
		}
	}
}
