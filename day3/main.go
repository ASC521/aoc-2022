package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
--- Day 3: Rucksack Reorganization ---

One Elf has the important job of loading all of the rucksacks with supplies for the jungle journey. Unfortunately, that Elf didn't quite follow the packing instructions, and so a few items now need to be rearranged.

Each rucksack has two large compartments. All items of a given type are meant to go into exactly one of the two compartments. The Elf that did the packing failed to follow this rule for exactly one item type per rucksack.

The Elves have made a list of all of the items currently in each rucksack (your puzzle input), but they need your help finding the errors. Every item type is identified by a single lowercase or uppercase letter (that is, a and A refer to different types of items).

The list of items for each rucksack is given as characters all on a single line. A given rucksack always has the same number of items in each of its two compartments, so the first half of the characters represent items in the first compartment, while the second half of the characters represent items in the second compartment.

For example, suppose you have the following list of contents from six rucksacks:

vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw

    The first rucksack contains the items vJrwpWtwJgWrhcsFMMfFFhFp, which means its first compartment contains the items vJrwpWtwJgWr, while the second compartment contains the items hcsFMMfFFhFp. The only item type that appears in both compartments is lowercase p.
    The second rucksack's compartments contain jqHRNqRjqzjGDLGL and rsFMfFZSrLrFZsSL. The only item type that appears in both compartments is uppercase L.
    The third rucksack's compartments contain PmmdzqPrV and vPwwTWBwg; the only common item type is uppercase P.
    The fourth rucksack's compartments only share item type v.
    The fifth rucksack's compartments only share item type t.
    The sixth rucksack's compartments only share item type s.

To help prioritize item rearrangement, every item type can be converted to a priority:

    Lowercase item types a through z have priorities 1 through 26.
    Uppercase item types A through Z have priorities 27 through 52.

In the above example, the priority of the item type that appears in both compartments of each rucksack is 16 (p), 38 (L), 42 (P), 22 (v), 20 (t), and 19 (s); the sum of these is 157.

Find the item type that appears in both compartments of each rucksack. What is the sum of the priorities of those item types?

*/

func getPriorities() map[string]int {
	priorities := map[string]int{}
	priority := 1
	for ch := 'a'; ch <= 'z'; ch++ {
		priorities[string(ch)] = priority
		priority += 1
	}
	for ch := 'a'; ch <= 'z'; ch++ {
		priorities[strings.ToUpper(string(ch))] = priority
		priority += 1
	}
	return priorities
}

func findMisplacedSupply(one, two string) string {

	supplies := make(map[string]bool)
	for _, s := range strings.Split(one, "") {
		supplies[s] = true
	}

	for _, s := range strings.Split(two, "") {
		if supplies[s] {
			return s
		}
	}
	return ""
}

func identifyBadgeGroup(one, two, three string) string {

	suppliesElfOne := make(map[string]bool)
	for _, s := range strings.Split(one, "") {
		suppliesElfOne[s] = true
	}

	suppliesElfTwo := make(map[string]bool)
	for _, s := range strings.Split(two, "") {
		suppliesElfTwo[s] = true
	}

	for _, s := range strings.Split(three, "") {
		if suppliesElfOne[s] && suppliesElfTwo[s] {
			return s
		}
	}

	panic("I should never get here.  It was guaranteed exactly 1 item would be duplicate in a group.")
}

func main() {

	part := 2
	priority := getPriorities()

	file, err := os.Open("day3/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rucksack := 1
	totalP := 0
	scanner := bufio.NewScanner(file)
	elfOne := ""
	elfTwo := ""
	elfThree := ""
	for scanner.Scan() {
		if part == 1 {
			fullInventory := scanner.Text()
			compartmentOne := fullInventory[0 : len(fullInventory)/2]
			compartmentTwo := fullInventory[len(fullInventory)/2:]
			misplacedSupply := findMisplacedSupply(compartmentOne, compartmentTwo)
			p := priority[misplacedSupply]
			totalP += p
			fmt.Printf("Rucksack %v: %v  %v\n", rucksack, misplacedSupply, p)
			rucksack += 1
		} else if part == 2 {
			line := scanner.Text()

			if elfOne == "" {
				elfOne = line
				continue
			} else if elfOne != "" && elfTwo == "" {
				elfTwo = line
				continue
			} else if elfOne != "" && elfTwo != "" && elfThree == "" {
				elfThree = line
				s := identifyBadgeGroup(elfOne, elfTwo, elfThree)
				totalP += priority[s]

				elfOne = ""
				elfTwo = ""
				elfThree = ""
			}

		}

	}
	fmt.Printf("Total Priority: %v\n", totalP)
}
