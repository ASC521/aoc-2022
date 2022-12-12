package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ClockCircuit struct {
	Cycle int
}

func NewClockCircuit() *ClockCircuit {
	return &ClockCircuit{
		Cycle: 1,
	}
}

func (c *ClockCircuit) CompleteCycle() {
	c.Cycle++
}

type CPU struct {
	X int
}

func NewCPU() *CPU {
	return &CPU{
		X: 1,
	}
}

func (c *CPU) Noop() {
}

func (c *CPU) AddX(amount int) {
	c.X += amount
}

type CRT struct {
	width  int
	height int
	pixels int
}

func NewCRT(width, height int) *CRT {
	return &CRT{width, height, width * height}
}

func (crt *CRT) Draw(cycle int, spriteLocation int) {
	leftPixel := spriteLocation - 1
	rightPixel := spriteLocation + 1
	pixelLoc := (cycle - 1) % crt.width
	if leftPixel <= pixelLoc && pixelLoc <= rightPixel {
		fmt.Print("#")
	} else {
		fmt.Print(".")
	}

	if cycle%40 == 0 {
		fmt.Print("\n")
	}
}

type VideoSystem struct {
	Circuit           *ClockCircuit
	CPU               *CPU
	CRT               *CRT
	SignalStrengthLog map[int]int
}

func NewVideoSystem(screenWidth, screenHeight int) *VideoSystem {
	return &VideoSystem{
		Circuit:           NewClockCircuit(),
		CPU:               NewCPU(),
		CRT:               NewCRT(screenWidth, screenHeight),
		SignalStrengthLog: map[int]int{},
	}
}

func (v *VideoSystem) RunProgram(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Split(line, " ")

		switch args[0] {
		case "addx":
			amount, err := strconv.Atoi(args[1])
			if err != nil {
				panic(err)
			}

			v.CRT.Draw(v.Circuit.Cycle, v.CPU.X)
			v.LogSignalStrength()
			v.Circuit.CompleteCycle()

			v.LogSignalStrength()
			v.CRT.Draw(v.Circuit.Cycle, v.CPU.X)
			v.CPU.AddX(amount)
			v.Circuit.CompleteCycle()
		case "noop":
			v.CPU.Noop()

			v.CRT.Draw(v.Circuit.Cycle, v.CPU.X)
			v.LogSignalStrength()
			v.Circuit.CompleteCycle()
		default:
			panic("I should never get here")
		}
	}
}

func (v *VideoSystem) LogSignalStrength() {
	v.SignalStrengthLog[v.Circuit.Cycle] = v.Circuit.Cycle * v.CPU.X
}

func main() {
	file, err := os.Open("/Users/alex.curto/code/aoc-2022/day10/input.txt")
	if err != nil {
		panic(err)
	}

	vs := NewVideoSystem(40, 6)
	vs.RunProgram(file)
	// fmt.Println(vs.SignalStrengthLog)
	// cycles := []int{20, 60, 100, 140, 180, 220}
	// sum := 0
	// for _, c := range cycles {
	// 	sum += vs.SignalStrengthLog[c]
	// }
	// fmt.Println(sum)
}
