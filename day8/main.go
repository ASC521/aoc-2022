package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
--- Day 8: Treetop Tree House ---

The expedition comes across a peculiar patch of tall trees all planted carefully in a grid. The Elves explain that a previous expedition planted these trees as a reforestation effort. Now, they're curious if this would be a good location for a tree house.

First, determine whether there is enough tree cover here to keep a tree house hidden. To do this, you need to count the number of trees that are visible from outside the grid when looking directly along a row or column.

The Elves have already launched a quadcopter to generate a map with the height of each tree (your puzzle input). For example:

30373
25512
65332
33549
35390

Each tree is represented as a single digit whose value is its height, where 0 is the shortest and 9 is the tallest.

A tree is visible if all of the other trees between it and an edge of the grid are shorter than it. Only consider trees in the same row or column; that is, only look up, down, left, or right from any given tree.

All of the trees around the edge of the grid are visible - since they are already on the edge, there are no trees to block the view. In this example, that only leaves the interior nine trees to consider:

    The top-left 5 is visible from the left and top. (It isn't visible from the right or bottom since other trees of height 5 are in the way.)
    The top-middle 5 is visible from the top and right.
    The top-right 1 is not visible from any direction; for it to be visible, there would need to only be trees of height 0 between it and an edge.
    The left-middle 5 is visible, but only from the right.
    The center 3 is not visible from any direction; for it to be visible, there would need to be only trees of at most height 2 between it and an edge.
    The right-middle 3 is visible from the right.
    In the bottom row, the middle 5 is visible, but the 3 and 4 are not.

With 16 trees visible on the edge and another 5 visible in the interior, a total of 21 trees are visible in this arrangement.

Consider your map; how many trees are visible from outside the grid?

*/

type Tree struct {
	height int
}

func NewTree(h int) *Tree {
	return &Tree{height: h}
}

type Forest struct {
	trees [][]*Tree
}

func NewForest() *Forest {
	return &Forest{
		trees: [][]*Tree{},
	}
}

func (f *Forest) AddTree(height, row int) {

	if len(f.trees) < row+1 {
		f.trees = append(f.trees, []*Tree{})
	}

	f.trees[row] = append(f.trees[row], NewTree(height))
}

func (f *Forest) String() string {
	str := ""
	for _, r := range f.trees {
		for _, t := range r {
			str += strconv.Itoa(t.height)
		}
		str += "\n"
	}

	return str
}

func (f *Forest) Visible(row, col int) bool {

	t := f.trees[row][col]
	if row == 0 || row == len(f.trees[row])-1 || col == 0 || col == len(f.trees)-1 {
		return true
	}

	top := 0
	for i := row - 1; i >= 0; i-- {
		if top < f.trees[i][col].height {
			top = f.trees[i][col].height
		}
	}

	down := 0
	for i := row + 1; i <= len(f.trees)-1; i++ {
		if down < f.trees[i][col].height {
			down = f.trees[i][col].height
		}
	}

	right := 0
	for i := col - 1; i >= 0; i-- {
		if right < f.trees[row][i].height {
			right = f.trees[row][i].height
		}
	}

	left := 0
	for i := col + 1; i <= len(f.trees[row])-1; i++ {
		if left < f.trees[row][i].height {
			left = f.trees[row][i].height
		}
	}

	if top < t.height || down < t.height || right < t.height || left < t.height {
		return true
	} else {
		return false
	}

}

func (f *Forest) TreesVisibleFromOutside() []*Tree {

	visibleTrees := []*Tree{}
	for row := 0; row <= len(f.trees)-1; row++ {
		for col := 0; col <= len(f.trees[row])-1; col++ {
			t := f.trees[row][col]
			if f.Visible(row, col) {
				visibleTrees = append(visibleTrees, t)
			}

		}

	}
	return visibleTrees
}

func (f *Forest) ScenicScore(row, col int) int {
	t := f.trees[row][col]
	if row == 0 || row == len(f.trees[row])-1 || col == 0 || col == len(f.trees)-1 {
		return 0
	}

	top := 0
	for i := row - 1; i >= 0; i-- {
		top++
		if t.height <= f.trees[i][col].height {
			break
		}
	}

	down := 0
	for i := row + 1; i <= len(f.trees)-1; i++ {
		down++
		if t.height <= f.trees[i][col].height {
			break
		}

	}

	left := 0
	for i := col - 1; i >= 0; i-- {
		left++
		if t.height <= f.trees[row][i].height {
			break
		}

	}

	right := 0
	for i := col + 1; i <= len(f.trees[row])-1; i++ {
		right++
		if t.height <= f.trees[row][i].height {
			break
		}
	}

	return top * right * down * left
}

func (f *Forest) ComputeScenicScores() []int {
	scores := []int{}
	for row := 0; row <= len(f.trees)-1; row++ {
		for col := 0; col <= len(f.trees[row])-1; col++ {
			score := f.ScenicScore(row, col)
			scores = append(scores, score)
		}

	}
	sort.Sort(sort.Reverse(sort.IntSlice(scores)))
	return scores
}

func main() {

	file, err := os.Open("/Users/alex.curto/code/aoc-2022/day8/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	forest := NewForest()
	scanner := bufio.NewScanner(file)
	rowNum := 0
	for scanner.Scan() {
		row := scanner.Text()
		nodes := strings.Split(row, "")
		for _, t := range nodes {
			h, err := strconv.Atoi(t)
			if err != nil {
				panic(err)
			}
			forest.AddTree(h, rowNum)
		}

		rowNum++

	}

	trees := forest.TreesVisibleFromOutside()
	fmt.Println(len(trees))
	scores := forest.ComputeScenicScores()
	fmt.Println(scores[0])

}
