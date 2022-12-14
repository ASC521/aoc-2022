package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Operation struct {
	op     string
	amount int
}

type Monkey struct {
	id             string
	holding        []int
	numInpsections int
	test           int
	operation      Operation
	true           int
	false          int
}

func NewMonkey(id string) *Monkey {
	return &Monkey{id: id}
}

func main() {

	file, err := os.Open("/Users/alex.curto/code/aoc-2022/day11/input.txt")
	if err != nil {
		panic(err)
	}

	monkeys := []*Monkey{}

	sc := bufio.NewScanner(file)
	var monkey *Monkey
	for sc.Scan() {
		line := sc.Text()

		switch {
		case strings.Contains(line, "Monkey "):
			id := strings.Split(line, " ")[1]

			monkey = NewMonkey(id)
			monkeys = append(monkeys, monkey)

		case strings.Contains(line, "Starting items:"):
			itemsStr := strings.Split(line, ": ")[1]
			itemsSplit := strings.Split(itemsStr, ", ")
			var items []int
			for _, i := range itemsSplit {
				item, err := strconv.Atoi(i)
				if err != nil {
					panic(err)
				}
				items = append(items, item)
			}

			monkey.holding = items

		case strings.Contains(line, "Operation:"):
			formula := strings.Split(line, "  Operation: ")[1]
			components := strings.Split(formula, " ")
			op := components[3]
			var amt int
			if components[4] == "old" {
				amt = 2
				op = "^"
			} else {
				amt, err = strconv.Atoi(components[4])
				if err != nil {
					panic(err)
				}
			}

			monkey.operation = Operation{op: op, amount: amt}

		case strings.Contains(line, "Test:"):
			divisor, err := strconv.Atoi(strings.Split(line, " ")[5])
			if err != nil {
				panic(err)
			}
			monkey.test = divisor

		case strings.Contains(line, "throw to monkey "):
			leftRight := strings.Split(line, ": ")
			trueFalse := leftRight[0][7:]
			idStr := strings.Split(leftRight[1], " ")[3]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				panic(err)
			}
			switch trueFalse {
			case "true":
				monkey.true = id
			case "false":
				monkey.false = id
			}

		case line == "":
			monkey = nil
		}

	}

	gcd := 1
	for _, monkey := range monkeys {
		gcd *= monkey.test
	}

	for i := 0; i < 10000; i++ {
		for _, m := range monkeys {
			for _, item := range m.holding {
				var wl int
				switch m.operation.op {
				case "*":
					wl = item * m.operation.amount
				case "+":
					wl = item + m.operation.amount
				case "^":
					if m.operation.amount == 2 {
						wl = item * item
					} else {
						wl = int(math.Pow(float64(item), float64(m.operation.amount)))
					}

				default:
					msg := fmt.Sprintf("Missing support for operation %v", m.operation.op)
					panic(msg)
				}

				// wl = wl / 3
				wl %= gcd

				var throwTo int
				if wl%m.test == 0 {
					throwTo = m.true
				} else {
					throwTo = m.false
				}

				monkeys[throwTo].holding = append(monkeys[throwTo].holding, wl)
				m.numInpsections++

			}
			m.holding = []int{}
		}
	}

	for _, m := range monkeys {
		fmt.Printf("Monkey %v : %v\n", m.id, m.holding)
	}

	var inspections []int
	for _, m := range monkeys {
		inspections = append(inspections, m.numInpsections)
		fmt.Printf("Monkey %v %v inspections\n", m.id, m.numInpsections)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(inspections)))
	fmt.Println(inspections)
	fmt.Println(inspections[0] * inspections[1])

}
