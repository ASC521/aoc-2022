package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
--- Day 4: Camp Cleanup ---

Space needs to be cleared before the last supplies can be unloaded from the ships, and so several Elves have been assigned the job of cleaning up sections of the camp. Every section has a unique ID number, and each Elf is assigned a range of section IDs.

However, as some of the Elves compare their section assignments with each other, they've noticed that many of the assignments overlap. To try to quickly find overlaps and reduce duplicated effort, the Elves pair up and make a big list of the section assignments for each pair (your puzzle input).

For example, consider the following list of section assignment pairs:

2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8

For the first few pairs, this list means:

    Within the first pair of Elves, the first Elf was assigned sections 2-4 (sections 2, 3, and 4), while the second Elf was assigned sections 6-8 (sections 6, 7, 8).
    The Elves in the second pair were each assigned two sections.
    The Elves in the third pair were each assigned three sections: one got sections 5, 6, and 7, while the other also got 7, plus 8 and 9.

This example list uses single-digit section IDs to make it easier to draw; your actual list might contain larger numbers. Visually, these pairs of section assignments look like this:

.234.....  2-4
.....678.  6-8

.23......  2-3
...45....  4-5

....567..  5-7
......789  7-9

.2345678.  2-8
..34567..  3-7

.....6...  6-6
...456...  4-6

.23456...  2-6
...45678.  4-8

Some of the pairs have noticed that one of their assignments fully contains the other. For example, 2-8 fully contains 3-7, and 6-6 is fully contained by 4-6. In pairs where one assignment fully contains the other, one Elf in the pair would be exclusively cleaning sections their partner will already be cleaning, so these seem like the most in need of reconsideration. In this example, there are 2 such pairs.

In how many assignment pairs does one range fully contain the other?

https://adventofcode.com/2022/day/4

*/

func parseLargeDigit(r string) map[int]bool {
	elf := make(map[int]bool)
	for _, a := range strings.Split(strings.Trim(r, "."), "") {
		i, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		elf[i] = true
	}

	return elf
}

func parseSignleDigit(r string) map[int]bool {
	elf := make(map[int]bool)
	lowHigh := strings.Split(r, "-")
	low, err := strconv.Atoi(lowHigh[0])
	if err != nil {
		panic(err)
	}

	high, err := strconv.Atoi(lowHigh[1])
	if err != nil {
		panic(err)
	}

	for i := low; i <= high; i++ {
		elf[i] = true
	}

	return elf
}

func contains(one, two map[int]bool) bool {

	for a := range one {
		if !two[a] {
			return false
		}
	}
	return true
}

func main() {

	file, err := os.Open("day4/large_section_assignments.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var assignmentOne string
	var assignmentTwo string
	fullyOverlappped := 0
	for scanner.Scan() {
		line := scanner.Text()
		assignments := strings.Split(line, ",")
		if len(assignments) == 2 {
			elfOne := parseSignleDigit(assignments[0])
			elfTwo := parseSignleDigit(assignments[1])

			if contains(elfOne, elfTwo) || contains(elfTwo, elfOne) {
				fmt.Printf("Assingments fully overlap: %v\n", line)
				fullyOverlappped += 1
			}

		} else if len(assignments) == 1 {
			if line == "" {
				assignmentOne = ""
				assignmentTwo = ""
				continue
			} else if assignmentOne != "" && assignmentTwo == "" {
				assignmentTwo = line
				elfOne := parseLargeDigit(assignmentOne)
				elfTwo := parseLargeDigit(assignmentTwo)
				if contains(elfOne, elfTwo) || contains(elfTwo, elfOne) {
					fmt.Printf("Assingments fully overlap: %v    %v\n", assignmentOne, assignmentTwo)
					fullyOverlappped += 1
				}
			} else if assignmentOne == "" && assignmentTwo == "" {
				assignmentOne = line
			}

		}
	}

	fmt.Printf("Fully Overlapping Assingments: %v\n", fullyOverlappped)

}
