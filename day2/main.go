package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Ripped off the boilerplate from ../day1 and made some tweaks.

type dialDirection int

const (
	dialDirectionInvalid dialDirection = iota
	dialDirectionLeft
	dialDirectionRight
)

func (d dialDirection) String() string {
	return [...]string{"DIAL_DIRECTION_INVALID", "DIAL_DIRECTION_LEFT", "DIAL_DIRECTION_RIGHT"}[d]
}

// L150, R1, L33, etc
func decodeLine(line string) (dialDirection, int, error) {
	line = strings.TrimSuffix(line, "\n")

	directionStr := line[0:1]
	direction := dialDirectionInvalid
	switch directionStr {
	case "L":
		direction = dialDirectionLeft
	case "R":
		direction = dialDirectionRight
	default:
		return dialDirectionInvalid, 0, fmt.Errorf("invalid direction %s", directionStr)
	}

	number, err := strconv.Atoi(line[1:])
	if err != nil {
		return dialDirectionInvalid, 0, fmt.Errorf("converting line string to number %q", err)
	}

	return direction, number, nil
}

// TODO: Could be fun to implement this under a DialSolver.Solve interface
// and implement other solutions, like using an InverseDial solution (L10 = R90 for dialMax=100)
// to see how the solutions compare in time complexity. For now, we only implement my first solution
// which passed the online checker.
func mainE(fileName string, dialDefault int, dialMax int, debug bool) error {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("reading file %q", err)
	}

	lines := strings.Lines(string(content))
	dial := dialDefault
	numZeros := 0
	for line := range lines {
		dialDirection, number, err := decodeLine(line)
		if err != nil {
			return fmt.Errorf("decoding line %q", err)
		}

		// Rotating by n results in
		// 1. crossing zero n times
		// 2. the same ending number as rotating by n % dialMax in either direction.
		numZeros += number / dialMax
		number %= dialMax
		oldDial := dial

		switch dialDirection {
		case dialDirectionRight:
			dial += number

			// Can still cross the threshold at most once
			// (e.g. for dialmax 100 and dialStart 50, 50 + 60 = 110 =~= 10),
			// Therefore, modulo by dialMax once more.
			if dial > dialMax {
				dial %= dialMax
				numZeros += 1
			}

		case dialDirectionLeft:
			dial -= number

			// If the counterclockwise dialing results in crossing the threshold (0),
			// then just add dialMax to it to invert it back into a valid dial position.
			// e.g. dialMax=100 dialStart=50: 50 -> L160 -> L60 -> -10, -10 + 100 = 90
			if dial < 0 {
				dial += dialMax
				numZeros += 1
			}
		default:
			return fmt.Errorf("unexpected direction %s", dialDirection)
		}

		if debug {
			fmt.Printf("direction=%s,number=%d,oldDial=%d,dial=%d,numZeros=%d\n", dialDirection, number, oldDial, dial, numZeros)
		}
	}

	fmt.Printf("The dial passed 0 %d times\n", numZeros)
	return nil
}

func main() {
	var dialMax, dialStart int
	var fileName string
	var debug bool

	flag.IntVar(&dialMax, "dial-max", 100, "the number of valid dial markers starting from [0 ex) 100 means [0-99] are valid")
	flag.IntVar(&dialStart, "dial-start", 50, "the dial marker to start the operation at")
	flag.StringVar(&fileName, "file-name", "input.txt", "name of the file that contains the combination list, must be in working directory")
	flag.BoolVar(&debug, "debug", false, "print out the dial before/after of each combo applied")

	flag.Parse()

	if err := mainE(fileName, dialStart, dialMax, debug); err != nil {
		log.Fatalf("failed to run %q", err)
	}

	fmt.Println("success!")
}
