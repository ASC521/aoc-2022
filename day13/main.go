package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func cmp(left, right any) int {

	ls, lOK := left.([]any)
	rs, rOK := right.([]any)

	switch {
	// Both INTS
	case !lOK && !rOK:
		return int(left.(float64) - right.(float64))
	// Left is an int
	case !lOK:
		ls = []any{left}

	// Right is an int
	case !rOK:
		rs = []any{right}
	}

	for i := 0; i < len(ls) && i < len(rs); i++ {
		if c := cmp(ls[i], rs[i]); c != 0 {
			return c
		}
	}

	return len(ls) - len(rs)
}

func main() {

	input, err := os.ReadFile("/Users/alex.curto/code/aoc-2022/day13/input.txt")
	if err != nil {
		log.Fatal(err)
	}

	pkts := []any{}
	part1 := 0
	for i, s := range strings.Split(strings.TrimSpace(string(input)), "\n\n") {
		s := strings.Split(s, "\n")
		var left, right any
		json.Unmarshal([]byte(s[0]), &left)
		json.Unmarshal([]byte(s[1]), &right)
		pkts = append(pkts, left, right)
		if cmp(left, right) <= 0 {
			part1 += i + 1
		}
	}

	fmt.Println(part1)
	pkts = append(pkts, []any{[]any{2.}}, []any{[]any{6.}})
	sort.Slice(pkts, func(i, j int) bool { return cmp(pkts[i], pkts[j]) < 0 })

	part2 := 1
	for i, p := range pkts {
		if fmt.Sprint(p) == "[[2]]" || fmt.Sprint(p) == "[[6]]" {
			part2 *= i + 1
		}
	}

	fmt.Println(part2)

}
