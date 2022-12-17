package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Square struct {
	elevation int
	row       int
	col       int
}

type Grid [][]*Square

func NewGrid() Grid {
	return [][]*Square{}
}

func (g Grid) FindNeighbors(s *Square) []*Square {

	n := []*Square{}

	if s.row < len(g)-1 {
		n = append(n, g[s.row+1][s.col])
	}

	if s.row > 0 {
		n = append(n, g[s.row-1][s.col])
	}

	if s.col < len(g[0])-1 {
		n = append(n, g[s.row][s.col+1])
	}

	if s.col > 0 {
		n = append(n, g[s.row][s.col-1])
	}

	return n

}

func elevationMap() map[string]int {
	elevations := map[string]int{}
	e := 1
	for ch := 'a'; ch <= 'z'; ch++ {
		elevations[string(ch)] = e
		e += 1
	}
	return elevations
}

func (g Grid) computeShortestDistance(start, end *Square) int {

	visited := map[*Square]bool{}
	unvisited := []*Square{start}

	distances := [][]int{}
	for r := 0; r <= len(g)-1; r++ {
		distances = append(distances, []int{})
		for c := 0; c <= len(g[r])-1; c++ {
			distances[r] = append(distances[r], -1)
		}
	}

	distances[start.row][start.col] = 0

	for {

		cs := unvisited[0]
		visited[cs] = true

		neighbors := g.FindNeighbors(cs)
		for _, n := range neighbors {

			if n.elevation > cs.elevation+1 {
				continue
			}

			if visited[n] || distances[n.row][n.col] != -1 {
				continue
			}

			unvisited = append(unvisited, n)

			distances[n.row][n.col] = distances[cs.row][cs.col] + 1

		}

		if cs == end {
			break
		}

		unvisited = unvisited[1:]

		if len(unvisited) == 0 {
			break
		}

	}

	return distances[end.row][end.col]

}

func main() {

	file, err := os.Open("/Users/alex.curto/code/aoc-2022/day12/input.txt")
	if err != nil {
		log.Fatal(err)
	}

	elevations := elevationMap()

	var starts []*Square
	var end *Square
	grid := NewGrid()
	sc := bufio.NewScanner(file)
	row := 0
	for sc.Scan() {
		line := sc.Text()
		col := 0
		grid = append(grid, []*Square{})
		for _, s := range strings.Split(line, "") {
			switch s {
			case "S":
				grid[row] = append(grid[row], &Square{
					elevation: elevations["a"],
					row:       row,
					col:       col,
				})
				starts = append(starts, grid[row][col])
			case "a":
				grid[row] = append(grid[row], &Square{
					elevation: elevations["a"],
					row:       row,
					col:       col,
				})
				starts = append(starts, grid[row][col])
			case "E":
				grid[row] = append(grid[row], &Square{
					elevation: elevations["z"],
					row:       row,
					col:       col,
				})
				end = grid[row][col]
			default:
				grid[row] = append(grid[row], &Square{
					elevation: elevations[s],
					row:       row,
					col:       col,
				})
			}

			col++

		}
		row++
	}

	distances := []int{}
	for _, start := range starts {
		d := grid.computeShortestDistance(start, end)
		distances = append(distances, d)
	}

	sort.Ints(distances)
	fmt.Println(distances)

}
