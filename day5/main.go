package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
--- Day 5: Supply Stacks ---

The expedition can depart as soon as the final supplies have been unloaded from the ships. Supplies are stored in stacks of marked crates, but because the needed supplies are buried under many other crates, the crates need to be rearranged.

The ship has a giant cargo crane capable of moving crates between stacks. To ensure none of the crates get crushed or fall over, the crane operator will rearrange them in a series of carefully-planned steps. After the crates are rearranged, the desired crates will be at the top of each stack.

The Elves don't want to interrupt the crane operator during this delicate procedure, but they forgot to ask her which crate will end up where, and they want to be ready to unload them as soon as possible so they can embark.

They do, however, have a drawing of the starting stacks of crates and the rearrangement procedure (your puzzle input). For example:

    [D]
[N] [C]
[Z] [M] [P]
 1   2   3

move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2

In this example, there are three stacks of crates. Stack 1 contains two crates: crate Z is on the bottom, and crate N is on top. Stack 2 contains three crates; from bottom to top, they are crates M, C, and D. Finally, stack 3 contains a single crate, P.

Then, the rearrangement procedure is given. In each step of the procedure, a quantity of crates is moved from one stack to a different stack. In the first step of the above rearrangement procedure, one crate is moved from stack 2 to stack 1, resulting in this configuration:

[D]
[N] [C]
[Z] [M] [P]
 1   2   3

In the second step, three crates are moved from stack 1 to stack 3. Crates are moved one at a time, so the first crate to be moved (D) ends up below the second and third crates:

        [Z]
        [N]
    [C] [D]
    [M] [P]
 1   2   3

Then, both crates are moved from stack 2 to stack 1. Again, because crates are moved one at a time, crate C ends up below crate M:

        [Z]
        [N]
[M]     [D]
[C]     [P]
 1   2   3

Finally, one crate is moved from stack 1 to stack 2:

        [Z]
        [N]
        [D]
[C] [M] [P]
 1   2   3

The Elves just need to know which crate will end up on top of each stack; in this example, the top crates are C in stack 1, M in stack 2, and Z in stack 3, so you should combine these together and give the Elves the message CMZ.

After the rearrangement procedure completes, what crate ends up on top of each stack?

https://adventofcode.com/2022/day/5
*/

var (
	ProcedureRX = regexp.MustCompile("move (?P<quantity>[0-9]{1,2}) from (?P<origStack>[0-9]{1,2}) to (?P<destStack>[0-9]{1,2})")
)

type stack struct {
	crates []string
	count  int
}

func NewStack() *stack {
	return &stack{}
}

func (s *stack) top() string {
	if len(s.crates) == 0 {
		return ""
	}
	return s.crates[0]
}

func (s *stack) push(crate string) {
	s.crates = append(s.crates, crate)
	s.count++
}

func (s *stack) prepend(crate []string) {
	s.crates = append(crate, s.crates...)
	s.count = s.count + len(crate)
}

func (s *stack) pop() string {
	if s.count == 0 {
		return ""
	}
	c := s.crates[0]
	s.crates = s.crates[1:s.count]
	s.count--
	return c
}

type warehouse struct {
	stacks []*stack
}

func NewWarehouse() *warehouse {
	return &warehouse{stacks: []*stack{}}
}

func (w *warehouse) String() string {
	var str string
	for i, s := range w.stacks {
		str += "Stack " + strconv.Itoa(i+1) + ": "
		for _, c := range s.crates {
			str += c + " "
		}

		str += "\n"
	}

	return str
}

func (w *warehouse) move(origStackLoc, destStackLoc, quantity int) {

	origStack := w.stacks[origStackLoc-1]
	destStack := w.stacks[destStackLoc-1]

	for i := 1; i <= quantity; i++ {
		c := origStack.pop()
		destStack.prepend([]string{c})
	}

}

func (w *warehouse) moveMultiple(origStackLoc, destStackLoc, quantity int) {
	origStack := w.stacks[origStackLoc-1]
	destStack := w.stacks[destStackLoc-1]

	crates := []string{}
	for i := 1; i <= quantity; i++ {
		crates = append(crates, origStack.pop())
	}
	destStack.prepend(crates)

}

func (w *warehouse) topMarks() string {

	var marks string
	for _, s := range w.stacks {
		marks += s.top()
	}
	return marks
}

func (w *warehouse) pushToStack(stack int, crate string) {

	if stack > len(w.stacks) {
		w.stacks = append(w.stacks, NewStack())
	}

	if crate != "" {
		w.stacks[stack-1].push(crate)
	}

}

type Procedure struct {
	quantity  int
	origStack int
	destStack int
}

func parseInput(file *os.File) (*warehouse, []Procedure) {
	w := NewWarehouse()
	p := []Procedure{}
	parseWarehouse := true
	r := bufio.NewReader(file)
	for {
		line, readerErr := r.ReadString('\n')
		if readerErr != nil && readerErr != io.EOF {
			panic(readerErr)
		}

		if line == "\n" {
			parseWarehouse = false
			continue
		}

		if strings.Contains(line, " 1   2   3 ") {
			continue
		}

		stack := 1
		if parseWarehouse {
			for i := 0; i <= len(line)-1; i = i + 4 {
				var crate string
				mark := line[i : i+3]
				if mark == "   " {
					crate = ""
				} else {
					crate = string(mark[1])
				}
				w.pushToStack(stack, crate)
				stack++
			}
		} else {

			if !strings.Contains(line, "move ") {
				continue
			}

			matches := ProcedureRX.FindStringSubmatch(line)
			quantity, err := strconv.Atoi(matches[1])
			if err != nil {
				panic(err)
			}

			origStack, err := strconv.Atoi(matches[2])
			if err != nil {
				panic(err)
			}

			destStack, err := strconv.Atoi(matches[3])
			if err != nil {
				panic(err)
			}

			p = append(p, Procedure{
				quantity:  quantity,
				origStack: origStack,
				destStack: destStack,
			})
		}

		if readerErr == io.EOF {
			break
		}
	}

	return w, p
}

func partOne(warehouse *warehouse, procedures []Procedure) {
	for _, p := range procedures {
		warehouse.move(p.origStack, p.destStack, p.quantity)
	}

	fmt.Println(warehouse)
	fmt.Println(warehouse.topMarks())
}

func partTwo(warehouse *warehouse, procedures []Procedure) {

	for _, p := range procedures {
		warehouse.moveMultiple(p.origStack, p.destStack, p.quantity)
	}

	fmt.Println(warehouse)
	fmt.Println(warehouse.topMarks())
}

func main() {
	part := 2
	file, err := os.Open("/Users/alex.curto/code/aoc-2022/day5/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	warehouse, procedures := parseInput(file)
	if part == 1 {
		partOne(warehouse, procedures)
	} else {
		partTwo(warehouse, procedures)
	}

}
