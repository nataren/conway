package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
	"flag"
)

func main() {
	var initialization string
	flag.StringVar(&initialization, "init", "random", "The initialization mode")

	var rows, columns int
	flag.IntVar(&rows, "rows", 20, "Number of rows")
	flag.IntVar(&columns, "columns", 20, "Number of rows")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	g := newGame(rows, columns)
	bootstrap(initialization, &g)

	fmt.Printf("\033[H\033[2J")

	fmt.Printf(g.String())
	time.Sleep(1000 * time.Millisecond)

	for {
		fmt.Printf(g.String())
		tick(&g)
		time.Sleep(50 * time.Millisecond)

		// Clear screen
		fmt.Printf("\033[H\033[2J")
	}

}

const dead = 0
const alive = 1

type game struct {
	rows int
	columns int
	board [][] byte
	tmpBoard [][] byte
}

func (g *game) String() string {
	var b strings.Builder
	for _, columns := range g.board {
		for _, state := range columns {
			var c rune
			if state == alive {
				c = ' '
			} else {
				c = '\u0240'
			}
			fmt.Fprintf(&b, "%c", c)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func newGame(rows, columns int) game {
	g := game{rows, columns, make([][]byte, rows), make([][]byte, rows)}
	for i := range g.board {
		g.board[i] = make([]byte, columns)
	}
	for i := range g.tmpBoard {
		g.tmpBoard[i] = make([]byte, columns)
	}
	return g
}

func tick(g *game) {
	for row_index, columns := range g.board {
		for column_index, state := range columns {
			g.tmpBoard[row_index][column_index] = nextState(row_index, column_index, state, g)
		}
	}
	outdatedBoard := g.board
	g.board = g.tmpBoard
	g.tmpBoard = outdatedBoard
}

func nextState(row, column int, currentState byte, g *game) byte {
	neighbors := getNeighbors(row, column, g)
	deadNeighbors := 0
	liveNeighbors := 0

	for _, neighbor := range neighbors {
		if neighbor.state == dead {
			deadNeighbors += 1
		} else if neighbor.state == alive {
			liveNeighbors +=1
		}
	}

	if currentState == dead {
		if liveNeighbors == 3 {
			return alive
		}
		return currentState
	} else if currentState == alive {
		if liveNeighbors < 2 {
			return dead
		}
		if liveNeighbors == 2 || liveNeighbors == 3 {
			return currentState
		}
		if liveNeighbors > 3 {
			return dead
		}
	}
	panic("wrong state")
}

type point struct {
	row int
	column int
}

type neighbor struct {
	state byte
	coordinate point
}

func (p *point) String() string {
	return fmt.Sprintf("(row=%d, column=%d)", p.column, p.row)
}

func getNeighbors(current_row, current_column int, g *game) []neighbor {
	var neighbors []neighbor
	for _, row := range [3]int{current_row - 1, current_row, current_row + 1} {
		for _, column := range [3]int{current_column - 1, current_column, current_column + 1} {
			// Skip itself
			if current_row == row && current_column == column {
				continue
			}
			// Check if inside board boundaries
			if row < 0 {
				continue
			}
			if column < 0 {
				continue
			}
			if row > g.rows - 1 {
				continue
			}
			if column > g.columns - 1 {
				continue
			}
			neighbors = append(neighbors, neighbor{g.board[row][column], point{row, column}})
		}
	}
	return neighbors
}

func bootstrap(method string, g *game) {
	rowOffset := g.rows / 2
	columnOffset := g.columns / 2

	if method == "diagonal" {
		for row, columns := range g.board {
			for column, _ := range columns {
				if row == column {
					g.board[row][column] = alive
				}
			}
		}
	} else if method == "random" {
		for row, columns := range g.board {
			for column, _ := range columns {
				g.board[row][column] = randomDeadOrAlive()
			}
		}
	} else if method == "glider" {
		g.board[0 + rowOffset][1 + columnOffset] = alive
		g.board[1 + rowOffset][2 + columnOffset] = alive
		g.board[2 + rowOffset][0 + columnOffset] = alive
		g.board[2 + rowOffset][1 + columnOffset] = alive
		g.board[2 + rowOffset][2 + columnOffset] = alive
	} else if method == "r-pentomino" {
		g.board[0 + rowOffset][1 + columnOffset] = alive
		g.board[0 + rowOffset][2 + columnOffset] = alive
		g.board[1 + rowOffset][0 + columnOffset] = alive
		g.board[1 + rowOffset][1 + columnOffset] = alive
		g.board[2 + rowOffset][1 + columnOffset] = alive
	} else if method == "diehard" {
		g.board[0 + rowOffset][6 + columnOffset] = alive
		g.board[1 + rowOffset][0 + columnOffset] = alive
		g.board[1 + rowOffset][1 + columnOffset] = alive
		g.board[2 + rowOffset][1 + columnOffset] = alive
		g.board[2 + rowOffset][5 + columnOffset] = alive
		g.board[2 + rowOffset][6 + columnOffset] = alive
		g.board[2 + rowOffset][7 + columnOffset] = alive
	} else if method == "acorn" {
		g.board[0 + rowOffset][1 + columnOffset] = alive
		g.board[1 + rowOffset][3 + columnOffset] = alive
		g.board[2 + rowOffset][0 + columnOffset] = alive
		g.board[2 + rowOffset][1 + columnOffset] = alive
		g.board[2 + rowOffset][4 + columnOffset] = alive
		g.board[2 + rowOffset][5 + columnOffset] = alive
		g.board[2 + rowOffset][6 + columnOffset] = alive
	} else if method == "gosper-glider-gun" {
		g.board[0][24] = alive

		g.board[1][22] = alive
		g.board[1][24] = alive

		g.board[2][12] = alive
		g.board[2][13] = alive
		g.board[2][20] = alive
		g.board[2][21] = alive
		g.board[2][34] = alive
		g.board[2][35] = alive

		g.board[3][11] = alive
		g.board[3][15] = alive
		g.board[3][20] = alive
		g.board[3][21] = alive
		g.board[3][34] = alive
		g.board[3][35] = alive

		g.board[4][0] = alive
		g.board[4][1] = alive
		g.board[4][10] = alive
		g.board[4][16] = alive
		g.board[4][20] = alive
		g.board[4][21] = alive

		g.board[5][0] = alive
		g.board[5][1] = alive
		g.board[5][10] = alive
		g.board[5][14] = alive
		g.board[5][16] = alive
		g.board[5][17] = alive
		g.board[5][22] = alive
		g.board[5][24] = alive

		g.board[6][10] = alive
		g.board[6][16] = alive
		g.board[6][24] = alive

		g.board[7][11] = alive
		g.board[7][15] = alive

		g.board[8][12] = alive
		g.board[8][13] = alive
	} else if method == "minimal-inf-1" {
		g.board[1 + rowOffset][7 + columnOffset] = alive
		g.board[2 + rowOffset][5 + columnOffset] = alive
		g.board[2 + rowOffset][7 + columnOffset] = alive
		g.board[2 + rowOffset][8 + columnOffset] = alive
		g.board[3 + rowOffset][5 + columnOffset] = alive
		g.board[3 + rowOffset][7 + columnOffset] = alive
		g.board[4 + rowOffset][5 + columnOffset] = alive
		g.board[5 + rowOffset][3 + columnOffset] = alive
		g.board[6 + rowOffset][1 + columnOffset] = alive
		g.board[6 + rowOffset][3 + columnOffset] = alive
	} else if method == "minimal-inf-2" {
		g.board[1 + rowOffset][1 + columnOffset] = alive
		g.board[1 + rowOffset][2 + columnOffset] = alive
		g.board[1 + rowOffset][3 + columnOffset] = alive
		g.board[1 + rowOffset][5 + columnOffset] = alive
		g.board[2 + rowOffset][1 + columnOffset] = alive
		g.board[3 + rowOffset][4 + columnOffset] = alive
		g.board[3 + rowOffset][5 + columnOffset] = alive
		g.board[4 + rowOffset][2 + columnOffset] = alive
		g.board[4 + rowOffset][3 + columnOffset] = alive
		g.board[4 + rowOffset][5 + columnOffset] = alive
		g.board[5 + rowOffset][1 + columnOffset] = alive
		g.board[5 + rowOffset][3 + columnOffset] = alive
		g.board[5 + rowOffset][5 + columnOffset] = alive
	} else if method == "minimal-inf-3" {
		g.board[1 + rowOffset][1 + columnOffset] = alive
		g.board[1 + rowOffset][2 + columnOffset] = alive
		g.board[1 + rowOffset][3 + columnOffset] = alive
		g.board[1 + rowOffset][4 + columnOffset] = alive
		g.board[1 + rowOffset][5 + columnOffset] = alive
		g.board[1 + rowOffset][6 + columnOffset] = alive
		g.board[1 + rowOffset][7 + columnOffset] = alive
		g.board[1 + rowOffset][8 + columnOffset] = alive
		g.board[1 + rowOffset][10 + columnOffset] = alive
		g.board[1 + rowOffset][11 + columnOffset] = alive
		g.board[1 + rowOffset][12 + columnOffset] = alive
		g.board[1 + rowOffset][13 + columnOffset] = alive
		g.board[1 + rowOffset][14 + columnOffset] = alive
		g.board[1 + rowOffset][18 + columnOffset] = alive
		g.board[1 + rowOffset][19 + columnOffset] = alive
		g.board[1 + rowOffset][20 + columnOffset] = alive
		g.board[1 + rowOffset][27 + columnOffset] = alive
		g.board[1 + rowOffset][28 + columnOffset] = alive
		g.board[1 + rowOffset][29 + columnOffset] = alive
		g.board[1 + rowOffset][30 + columnOffset] = alive
		g.board[1 + rowOffset][31 + columnOffset] = alive
		g.board[1 + rowOffset][32 + columnOffset] = alive
		g.board[1 + rowOffset][33 + columnOffset] = alive
		g.board[1 + rowOffset][35 + columnOffset] = alive
		g.board[1 + rowOffset][36 + columnOffset] = alive
		g.board[1 + rowOffset][37 + columnOffset] = alive
		g.board[1 + rowOffset][38 + columnOffset] = alive
		g.board[1 + rowOffset][39 + columnOffset] = alive
	}
}

func randomDeadOrAlive() byte {
	n := 100
	r := rand.Intn(n)
	if r < n / 2 {
		return alive
	}
	return dead
}