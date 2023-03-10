package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	filename string
)

const (
	left  = -1
	right = 1
)

type instruction struct {
	// state that has to been matched
	state string

	// symbol that has to been matched
	symbol string

	// move that must be made after execution
	move int

	// symbol that is going to be in cell after instruction execution
	newSymbol string

	// newState that is going to be after instruction execution
	newState string
}

func (i instruction) String() string {
	var move string
	switch i.move {
	case left:
		move = "<-"
	case right:
		move = "->"
	default:
		panic("Unexpected move value")
	}
	return fmt.Sprintf("%s %s %s %s %s", i.state, i.symbol, i.newSymbol, move, i.newState)
}

type machine struct {
	tape  []string
	head  uint
	state string
}

func (m machine) match(i instruction) bool {
	return m.state == i.state && m.tape[m.head] == i.symbol
}

func (m *machine) execute(i instruction) {
	if m.head == 0 && i.move < 0 {
		panic("tape underflow")
	}
	m.tape[m.head], m.state = i.newSymbol, i.newState
	m.head = uint(int(m.head) + i.move)
}

func (m machine) print() {
	fmt.Printf("%s | ", m.state)

	offset := len(m.state) + 3
	for i, s := range m.tape {
		fmt.Printf("%s ", s)
		if uint(i) < m.head {
			offset += len(s) + 1
		}
	}
	fmt.Println()

	fmt.Printf("%s%s\n", strings.Repeat(" ", offset), strings.Repeat("^", len(m.tape[m.head])))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "aurochs: error: missing instruction file parameter")
		os.Exit(1)
	}

	filename = os.Args[1]
	buf, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "aurochs: error: cannot read instruction file '%s': %v", filename, err)
		os.Exit(1)
	}

	instructions := []instruction{}
	var initialState string

	for lineno, line := range strings.Split(string(buf), "\n") {
		lineno += 1
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "//") {
			continue
		}

		tokens := strings.Split(line, " ")
		if len(tokens) != 5 {
			fmt.Fprintf(os.Stderr, "%s:%d: error: cannot match instruction: %s\n", filename, lineno, line)
			fmt.Fprintf(os.Stderr, "aurochs: note: instruction format is: STATE CELL CELL MOVE STATE\n")
			os.Exit(1)
		}

		var move int
		switch tokens[3] {
		case "<-":
			move = left
		case "->":
			move = right
		default:
			fmt.Fprintf(os.Stderr, "%s:%d: error: move not recognized: %s\n", filename, lineno, tokens[2])
			os.Exit(1)
		}

		instructions = append(instructions, instruction{
			state:     tokens[0],
			symbol:    tokens[1],
			move:      move,
			newSymbol: tokens[2],
			newState:  tokens[4],
		})

		if len(instructions) == 1 {
			initialState = instructions[0].state
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s> ", initialState)
		os.Stdout.Sync()
		if !scanner.Scan() {
			break
		}
		m := machine{
			tape:  strings.Split(strings.TrimSpace(scanner.Text()), " "),
			head:  0,
			state: initialState,
		}
	next:
		for {
			for _, instruction := range instructions {
				if m.match(instruction) {
					m.execute(instruction)
					m.print()
					continue next
				}
			}
			break
		}
	}
}
