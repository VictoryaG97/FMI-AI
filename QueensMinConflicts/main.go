package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

var (
	numberOfQueens = flag.Int("n", 8, "number of queens")
	maxIter int
	queensByRow []int
	queensByLeftDiagonals []int
	queensByRightDiagonals []int

)

func QueensByRowAndDiagonals(positions []int) ([]int, []int, []int) {
	queensOnRow := make([]int, *numberOfQueens)
	queensOnLeftDiagonals := make([]int, (*numberOfQueens * 2) - 1)
	queensOnRightDiagonals := make([]int, (*numberOfQueens * 2) - 1)

	for col, row := range positions {
		sumL := row + col
		sumR := *numberOfQueens - row - 1 + col

		queensOnRow[row] += 1
		queensOnLeftDiagonals[sumL] += 1
		queensOnRightDiagonals[sumR] += 1
	}

	return queensOnRow, queensOnLeftDiagonals, queensOnRightDiagonals
}

func findMaxConflicts(positions []int) int {
	max := 0
	var maxConflicts []int

	for col, row := range positions {
		currConflicts := 0
		sumL := row + col
		sumR := *numberOfQueens - row - 1 + col

		currConflicts += queensByRow[row] - 1
		currConflicts += queensByLeftDiagonals[sumL] - 1
		currConflicts += queensByRightDiagonals[sumR] - 1

		if max == currConflicts {
			maxConflicts = append(maxConflicts, col)
		} else if max < currConflicts {
			max = currConflicts
			maxConflicts = maxConflicts[:0]
			maxConflicts = append(maxConflicts, col)
		}
	}

	if max == 0 {
		return -1
	}

	rand.Seed(time.Now().UnixNano())
	return maxConflicts[rand.Intn(len(maxConflicts))]
}

func findMinConflicts(maxConflictsCol int) int {
	min := *numberOfQueens
	var minConflicts []int

	for i := 0; i < *numberOfQueens; i++ {
		currConflicts := 0
		sumL := i + maxConflictsCol
		sumR := *numberOfQueens - i - 1 + maxConflictsCol

		currConflicts += queensByRow[i]
		currConflicts += queensByLeftDiagonals[sumL]
		currConflicts += queensByRightDiagonals[sumR]

		if min == currConflicts {
			minConflicts = append(minConflicts, i)
		} else if min > currConflicts {
			min = currConflicts
			minConflicts = minConflicts[:0]
			minConflicts = append(minConflicts, i)
		}
	}

	rand.Seed(time.Now().UnixNano())
	return minConflicts[rand.Intn(len(minConflicts))]
}

func removeQueen(col int, positions []int) {
	row := positions[col]
	sumL := row + col
	sumR := *numberOfQueens - row - 1 + col

	queensByRow[row] -= 1
	queensByLeftDiagonals[sumL] -= 1
	queensByRightDiagonals[sumR] -= 1
}

func makeAMove(maxConflictsCol int, minConflictsRow int, positions []int) []int {
	newSumL  := maxConflictsCol + minConflictsRow
	newSumR := *numberOfQueens - minConflictsRow - 1 + maxConflictsCol

	positions[maxConflictsCol] = minConflictsRow

	queensByRow[minConflictsRow] += 1
	queensByLeftDiagonals[newSumL]  += 1
	queensByRightDiagonals[newSumR]  += 1

	return positions
}

func PrintBoard(positions []int) {
	for row := 0; row < *numberOfQueens; row++ {
		for col := 0; col < *numberOfQueens; col++ {
			if positions[col] == row {
				fmt.Print("* ")
			} else {
				fmt.Printf("_ ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func Play() bool {
	positions := make([]int, *numberOfQueens)

	for i := 0; i < *numberOfQueens; i++ {
		if i % 2 == 0 {
			positions[i] = *numberOfQueens / 2
		} else {
			positions[i] = (*numberOfQueens / 2) + 1
		}
	}

	queensByRow, queensByLeftDiagonals, queensByRightDiagonals = QueensByRowAndDiagonals(positions)
	startUpTime := time.Now()

	i := 0
	for i < maxIter {
		i++

		// get queen with max conflicts
		queenWithMaxConflicts := findMaxConflicts(positions)

		if queenWithMaxConflicts == -1 {
			curTime := time.Now().Sub(startUpTime).Milliseconds()
			PrintBoard(positions)
			logrus.Printf("Algorithm finished in %d miliseconds\n", curTime)
			return true
		}

		removeQueen(queenWithMaxConflicts, positions)

		// move on row with min conflicts
		moveToRow := findMinConflicts(queenWithMaxConflicts)
		positions = makeAMove(queenWithMaxConflicts, moveToRow, positions)
	}

	if findMaxConflicts(positions) == -1 {
		curTime := time.Now().Sub(startUpTime).Milliseconds()
		PrintBoard(positions)
		logrus.Printf("Algorithm finished in %d miliseconds\n", curTime)
		return true
	}
	return false
}

func main() {
	flag.Parse()

	maxIter = *numberOfQueens * 3

	for !Play() {
	}
}