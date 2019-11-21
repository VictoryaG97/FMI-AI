package main

import (
	"container/heap"
	"flag"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

var (
	numberedTails = flag.Int("n", 8, "count of numbered tails (like 8, 15, 24, 35, etc.)")
	zeroIndex     = flag.Int("i", -1, "index of the empty (0) tile")
	dimension int
	goalMap   BoardMap
	initialHeuristic int
)

type BoardMap map[int][2]int

type Board struct {
	board BoardMap
	allMoves []string
	parent string
	cost      int
	heuristic int
	f         int
	index	  int
}

//////////////////////////////
//    Priority Queue       //
/////////////////////////////

type PriorityQueue []*Board

func (pQ PriorityQueue) Len() int { return len(pQ) }

func (pQ PriorityQueue) Less(i, j int) bool {
	return pQ[i].f < pQ[j].f
}

func (pQ *PriorityQueue) Push(board interface{}) {
	n := len(*pQ)
	newBoard := board.(Board)
	newBoard.index = n
	*pQ = append(*pQ, &newBoard)
}

func (pQ *PriorityQueue) Pop() interface{} {
	old := *pQ
	n := len(old)
	oldBoard := old[n-1]
	old[n-1] = nil
	oldBoard.index = -1
	*pQ = old[0 : n-1]
	return oldBoard
}

func (pQ PriorityQueue) Swap(i, j int) {
	pQ[i], pQ[j] = pQ[j], pQ[i]
	pQ[i].index = i
	pQ[j].index = j
}

func (pQ *PriorityQueue) Update(b *Board, f int) {
	b.f = f
	heap.Fix(pQ, b.index)
}

func (b Board) getMoves() []string{
	var moves []string

	if b.board[0][0] > 0 && b.parent != "RIGHT"	     	  	 { moves = append(moves, "LEFT")}
	if b.board[0][0] < (dimension - 1) && b.parent != "LEFT" { moves = append(moves, "RIGHT")}
	if b.board[0][1] > 0 && b.parent != "DOWN"		     	 { moves = append(moves, "UP")}
	if b.board[0][1] < (dimension - 1) && b.parent != "UP"   { moves = append(moves, "DOWN")}

	return moves
}

func makeAMove(board BoardMap, move string) BoardMap{
	var tailToMove int
	newBoard := BoardMap{}

	curZeroCol := board[0][0]
	curZeroRow := board[0][1]
	zeroCol := curZeroCol
	zeroRow := curZeroRow

	if move == "LEFT" {
		zeroCol = curZeroCol - 1
	} else if move == "RIGHT" {
		zeroCol = curZeroCol + 1
	}else if move == "UP" {
		zeroRow = curZeroRow -1
	} else {
		zeroRow = curZeroRow + 1
	}

	for number, coords := range board {
		if coords[0] == zeroCol && coords[1] == zeroRow {
			tailToMove = number
		} else {
			newBoard[number] = coords
		}
	}

	newBoard[tailToMove] = [2]int{curZeroCol, curZeroRow}
	newBoard[0] = [2]int{zeroCol, zeroRow}

	return newBoard
}

func generateChildren(currBoard Board) []Board{
	var children []Board
	possibleMoves := currBoard.getMoves()
	for _, move := range possibleMoves {
		childBoard := currBoard
		childBoard.board = makeAMove(childBoard.board, move)
		if childBoard.board != nil{
			childBoard.cost += 1
			childBoard.parent = move
			childBoard.allMoves = append(currBoard.allMoves, move)

			childBoard.heuristic = calculateManhattanDistance(childBoard.board)
			childBoard.f = childBoard.cost + childBoard.heuristic
			if childBoard.f <= initialHeuristic {
				children = append(children, childBoard)
			}
		}
	}

	return children
}

func setIndexes(board []int) BoardMap {
	m := make(BoardMap)
	for i := 0; i < len(board); i++ {
		row := i / dimension
		col := i % dimension
		m[board[i]] = [2]int{col, row}
	}
	return m
}

func calculateManhattanDistance(inputMap BoardMap) int {
	var md int
	for i := 1; i < dimension*dimension; i++ {
		md += int(
			math.Abs(float64(inputMap[i][0]-goalMap[i][0])) +
				math.Abs(float64(inputMap[i][1]-goalMap[i][1])))
	}
	return  md
}

func Pop(boards *[]Board) interface{} {
	old := *boards
	oldBoard := old[0]
	*boards = old[1:]

	return oldBoard
}

func Push(boards *[]Board, board interface{}) {
	newBoard := board.(Board)
	*boards = append(*boards, newBoard)
}

func exists(oldStates []BoardMap, currState BoardMap) bool {
	for _, board := range oldStates {
		if cmp.Equal(board, currState) { return true}
	}
	return false
}

func main() {
	flag.Parse()

	dimension = int(math.Sqrt(float64(*numberedTails+1)))

	input := make([]int, *numberedTails + 1)

	fmt.Println("Enter the input board")
	for i := 0; i <= *numberedTails; i++ {
		_, err := fmt.Scan(&input[i])
		if err != nil {
			fmt.Printf("Couldn't get the input board %v\n", err)
			return
		}
	}

	goalBoard := make([]int, *numberedTails)
	for i := 0; i < *numberedTails; i++ {
		goalBoard[i] = i + 1
	}
	if *zeroIndex != -1 && *zeroIndex != *numberedTails {
		tempBoard := make([]int, *zeroIndex+1)
		copy(tempBoard, goalBoard[:*zeroIndex])
		goalBoard = append(tempBoard, goalBoard[*zeroIndex:]...)
	} else {
		goalBoard = append(goalBoard, 0)
	}

	startUpTime := time.Now()
	goalMap  = setIndexes(goalBoard)

	inputBoard := Board{
		cost: 0,
		board: setIndexes(input),
		parent: "",
	}
	inputBoard.heuristic = calculateManhattanDistance(inputBoard.board)
	inputBoard.f = inputBoard.cost + inputBoard.heuristic
	initialHeuristic = inputBoard.heuristic

	var oldStates []BoardMap
	goingStates := make([]Board, 1)
	currState := Board{}

	goingStates[0] = inputBoard

	for {
		if len(goingStates) > 0 {
			currState = Pop(&goingStates).(Board)
			if cmp.Equal(currState.board, goalMap) {
				break
			}

			for _, child := range generateChildren(currState) {
				if !exists(oldStates, currState.board) {
					Push(&goingStates, child)
				}
			}
			oldStates = append(oldStates, currState.board)
		} else {
			initialHeuristic += 1
			Push(&goingStates, inputBoard)
			oldStates = oldStates[:0]
		}
	}
	fmt.Printf("Cost:    %v\n", currState.cost)
	fmt.Printf("Moves:    %v\n", currState.allMoves)
	logrus.Printf("Algorithm finished in %d miliseconds\n", time.Now().Sub(startUpTime).Milliseconds())
}
