package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
--- Day 9: Rope Bridge ---

This rope bridge creaks as you walk along it. You aren't sure how old it is, or whether it can even support your weight.

It seems to support the Elves just fine, though. The bridge spans a gorge which was carved out by the massive river far below you.

You step carefully; as you do, the ropes stretch and twist. You decide to distract yourself by modeling rope physics; maybe you can even figure out where not to step.

Consider a rope with a knot at each end; these knots mark the head and the tail of the rope. If the head moves far enough away from the tail, the tail is pulled toward the head.

Due to nebulous reasoning involving Planck lengths, you should be able to model the positions of the knots on a two-dimensional grid. Then, by following a hypothetical series of motions (your puzzle input) for the head, you can determine how the tail will move.

Due to the aforementioned Planck lengths, the rope must be quite short; in fact, the head (H) and tail (T) must always be touching (diagonally adjacent and even overlapping both count as touching):

....
.TH.
....

....
.H..
..T.
....

...
.H. (H covers T)
...

If the head is ever two steps directly up, down, left, or right from the tail, the tail must also move one step in that direction so it remains close enough:

.....    .....    .....
.TH.. -> .T.H. -> ..TH.
.....    .....    .....

...    ...    ...
.T.    .T.    ...
.H. -> ... -> .T.
...    .H.    .H.
...    ...    ...

Otherwise, if the head and tail aren't touching and aren't in the same row or column, the tail always moves one step diagonally to keep up:

.....    .....    .....
.....    ..H..    ..H..
..H.. -> ..... -> ..T..
.T...    .T...    .....
.....    .....    .....

.....    .....    .....
.....    .....    .....
..H.. -> ...H. -> ..TH.
.T...    .T...    .....
.....    .....    .....

You just need to work out where the tail goes as the head follows a series of motions. Assume the head and the tail both start at the same position, overlapping.

For example:

R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2

This series of motions moves the head right four steps, then up four steps, then left three steps, then down one step, and so on. After each step, you'll need to update the position of the tail if the step means the head is no longer adjacent to the tail. Visually, these motions occur as follows (s marks the starting position as a reference point):

== Initial State ==

......
......
......
......
H.....  (H covers T, s)

== R 4 ==

......
......
......
......
TH....  (T covers s)

......
......
......
......
sTH...

......
......
......
......
s.TH..

......
......
......
......
s..TH.

== U 4 ==

......
......
......
....H.
s..T..

......
......
....H.
....T.
s.....

......
....H.
....T.
......
s.....

....H.
....T.
......
......
s.....

== L 3 ==

...H..
....T.
......
......
s.....

..HT..
......
......
......
s.....

.HT...
......
......
......
s.....

== D 1 ==

..T...
.H....
......
......
s.....

== R 4 ==

..T...
..H...
......
......
s.....

..T...
...H..
......
......
s.....

......
...TH.
......
......
s.....

......
....TH
......
......
s.....

== D 1 ==

......
....T.
.....H
......
s.....

== L 5 ==

......
....T.
....H.
......
s.....

......
....T.
...H..
......
s.....

......
......
..HT..
......
s.....

......
......
.HT...
......
s.....

......
......
HT....
......
s.....

== R 2 ==

......
......
.H....  (H covers T)
......
s.....

......
......
.TH...
......
s.....

After simulating the rope, you can count up all of the positions the tail visited at least once. In this diagram, s again marks the starting position (which the tail also visited) and # marks other positions the tail visited:

..##..
...##.
.####.
....#.
s###..

So, there are 13 positions the tail visited at least once.

Simulate your complete hypothetical series of motions. How many positions does the tail of the rope visit at least once?
*/

type Direction string

const (
	Right = Direction("right")
	Up    = Direction("up")
	Down  = Direction("down")
	Left  = Direction("left")
)

var (
	StayPut = []Distance{
		{x: 0, y: 0},
		{x: 1, y: 1},
		{x: 1, y: -1},
		{x: -1, y: -1},
		{x: -1, y: 1},
		{x: -1, y: 0},
		{x: 1, y: 0},
		{x: 0, y: 1},
		{x: 0, y: -1},
	}
)

type Position struct {
	x, y int
}

func (p Position) String() string {
	return fmt.Sprintf("{ x: %v, y: %v }", p.x, p.y)
}

type Distance struct {
	x, y int
}

func (d Distance) String() string {
	return fmt.Sprintf("{ x: %v, y: %v }", d.x, d.y)
}

func DistanceBetween(position, other Position) Distance {
	xDist := position.x - other.x
	yDist := position.y - other.y
	return Distance{
		x: xDist,
		y: yDist,
	}
}

type Planck struct {
	pos  Position
	tail *Planck
	name string
}

func NewPlanck(name string) *Planck {
	return &Planck{
		pos:  Position{0, 0},
		tail: nil,
		name: name,
	}
}

func (p *Planck) MoveRight() {
	p.pos.x += 1
}

func (p *Planck) MoveDown() {
	p.pos.y -= 1
}

func (p *Planck) MoveLeft() {
	p.pos.x -= 1
}

func (p *Planck) MoveUp() {
	p.pos.y += 1
}

func (p *Planck) MoveUpAndRight() {
	p.MoveUp()
	p.MoveRight()
}

func (p *Planck) MoveUpAndLeft() {
	p.MoveUp()
	p.MoveLeft()
}

func (p *Planck) MoveDownAndRight() {
	p.MoveDown()
	p.MoveRight()
}

func (p *Planck) MoveDownAndLeft() {
	p.MoveDown()
	p.MoveLeft()
}

type Bridge struct {
	head            *Planck
	tailPositionLog map[string]bool
}

func NewBridge() *Bridge {
	return &Bridge{
		head: nil,
		tailPositionLog: map[string]bool{
			Position{0, 0}.String(): true,
		},
	}
}

func (b *Bridge) AddPlanck(name string) {
	p := NewPlanck(name)
	if b.head == nil {
		b.head = p
	} else {
		parent := b.head
		for parent.tail != nil {
			parent = parent.tail
		}
		parent.tail = p
	}
}

func (b *Bridge) MoveHead(direction Direction, amount int) {
	move := fmt.Sprintf("****** %v %v ******\n", direction, amount)
	fmt.Printf(move)
	for i := 1; i <= amount; i++ {
		switch direction {
		case Right:
			b.head.MoveRight()
		case Left:
			b.head.MoveLeft()
		case Down:
			b.head.MoveDown()
		case Up:
			b.head.MoveUp()
		default:
			msg := fmt.Sprintf("I have exhausted all possible directions.  Should not have gotten here.  Direction recieved %v", direction)
			panic(msg)
		}

		h := b.head
		for h.tail != nil {
			if h.name == "9" {
				fmt.Println("Debug")
			}
			b.MoveTail(h)
			h = h.tail
		}
		b.tailPositionLog[h.pos.String()] = true
		b.Show()

	}
	fmt.Printf(move)

}

func (b *Bridge) MoveTail(head *Planck) {

	dis := DistanceBetween(head.pos, head.tail.pos)
	t := head.tail

	switch {
	case dis == Distance{x: 2, y: 0}:
		t.MoveRight()
	case dis == Distance{x: 2, y: 1}:
		t.MoveUpAndRight()
	case dis == Distance{x: 2, y: -1}:
		t.MoveDownAndRight()
	case dis == Distance{x: 2, y: 2}:
		t.MoveUpAndRight()

	case dis == Distance{x: -2, y: 0}:
		t.MoveLeft()
	case dis == Distance{x: -2, y: -1}:
		t.MoveDownAndLeft()
	case dis == Distance{x: -2, y: 1}:
		t.MoveUpAndLeft()
	case dis == Distance{x: -2, y: -2}:
		t.MoveDownAndLeft()

	case dis == Distance{x: 0, y: 2}:
		t.MoveUp()
	case dis == Distance{x: 1, y: 2}:
		t.MoveUpAndRight()
	case dis == Distance{x: -1, y: 2}:
		t.MoveUpAndLeft()
	case dis == Distance{x: -2, y: 2}:
		t.MoveUpAndLeft()

	case dis == Distance{x: 0, y: -2}:
		t.MoveDown()
	case dis == Distance{x: 1, y: -2}:
		t.MoveDownAndRight()
	case dis == Distance{x: -1, y: -2}:
		t.MoveDownAndLeft()
	case dis == Distance{x: 2, y: -2}:
		t.MoveDownAndRight()

	case ContainsDistance(StayPut, dis):

	default:
		msg := fmt.Sprintf("Got a distance I should never get.  Dis: %v  Head: %v  Tail: %v", dis, b.head.pos, b.head.tail.pos)
		panic(msg)
	}
}

func (b *Bridge) Show() {

	if b.head != nil {
		p := b.head
		for p.tail != nil {
			fmt.Printf("{name: %v {x: %v y: %v }} ", p.name, p.pos.x, p.pos.y)
			p = p.tail
		}
	}
	fmt.Println()
}

func ContainsDistance(distances []Distance, f Distance) bool {
	for _, d := range distances {
		if d == f {
			return true
		}
	}

	return false
}

func main() {

	bridge := NewBridge()
	file, err := os.Open("/Users/alex.curto/code/aoc-2022/day9/input.txt")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		var name string
		if i == 0 {
			name = "H"
		} else {
			name = strconv.Itoa(i)
		}
		bridge.AddPlanck(name)
	}

	bridge.Show()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Split(line, " ")
		var direction Direction
		switch args[0] {
		case "R":
			direction = Right
		case "U":
			direction = Up
		case "L":
			direction = Left
		case "D":
			direction = Down
		default:
			panic("Parsed a direction I do not support")
		}

		amount, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}

		bridge.MoveHead(direction, amount)
	}

	bridge.PrintHits()
	fmt.Println(len(bridge.tailPositionLog))

}

func (b *Bridge) PrintHits() {

	top := 14
	bottom := -5
	left := -11
	right := 14
	rowNum := top
	for y := top; y >= bottom; y-- {
		// fmt.Printf("%v  ", rowNum)
		rowNum--
		for x := left; x <= right; x++ {
			key := fmt.Sprintf("{ x: %v, y: %v }", x, y)

			if x == 0 && y == 0 {
				fmt.Print("s")
			} else if b.tailPositionLog[key] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
