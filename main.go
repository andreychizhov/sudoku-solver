package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Field [9][9]int

type Cell struct { // Cell
	X int
	Y int
}

type Area struct {
	Tl Cell
	Br Cell
}

func main() {
	var field Field
	field = Field{}

	data, err := ReadFile("input.txt")
	if err != nil {
		fmt.Print("Error reading from file...")
		os.Exit(1)
	}
	field.Fill(data)

	field.Print()
	fmt.Print(field.CanPutIntoCell(4, 4, 5))

}

func (f *Field) Solve() {
	areas := GetSquares()

	for _, a := range areas {

	}

}

func GetHints(f Field) map[Cell]int {
	var hints map[Cell]int
	hints = make(map[Cell]int)

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if f[i][j] > 0 {
				hints[Cell{i, j}] = f[i][j]
			}
		}
	}
	return hints
}

func (f Field) CanPutIntoCell(x, y, n int) bool {
	if f[x][y] > 0 {
		return false
	}

	var areas = GetAreas(Cell{x, y})

	for _, a := range areas {
		if ExistsInArea(f, a, n) {
			return false
		}
	}
	return true
}

// Check if number already exists in given area
func ExistsInArea(f Field, a Area, n int) bool {
	for i := a.Tl.X; i <= a.Br.X; i++ {
		for j := a.Tl.Y; j <= a.Br.Y; j++ {
			if f[i][j] == n {
				return true
			}
		}
	}
	return false
}

// Get row, column and square that contains given cell
func GetAreas(p Cell) []Area {
	var areas []Area

	areas = append(areas, Area{Cell{p.X, 0}, Cell{p.X, 8}})
	areas = append(areas, Area{Cell{0, p.Y}, Cell{8, p.Y}})

	for i := 0; i <= 6; i += 3 {
		for j := 0; j <= 6; j += 3 {
			if p.X >= i && p.X <= i+2 && p.Y >= j && p.Y <= j+2 {
				areas = append(areas, Area{Cell{i, j}, Cell{i + 2, j + 2}})
			}
		}
	}
	return areas
}

// Get all square areas
func GetSquares() []Area {
	areas := make([]Area, 9, 9)
	for i := 0; i <= 6; i += 3 {
		for j := 0; j <= 6; j += 3 {
			areas = append(areas, Area{Cell{i, j}, Cell{i + 2, j + 2}})
		}
	}
	return areas
}

func (f *Field) Check() (bool, Cell, int) {

	for i := 0; i < 9; i++ {
		if s := Sum(*f, Cell{i, 0}, Cell{i, 8}); s != 45 {
			return false, Cell{i, 8}, s
		}
		if s := Sum(*f, Cell{0, i}, Cell{8, i}); s != 45 {
			return false, Cell{8, i}, s
		}
	}

	for i := 0; i <= 6; i += 3 {
		for j := 0; j <= 6; j += 3 {
			if s := Sum(*f, Cell{i, j}, Cell{i + 2, j + 2}); s != 45 {
				return false, Cell{i, j}, s
			}
		}
	}

	return true, Cell{0, 0}, 45
}

func Sum(field Field, p1, p2 Cell) int {
	sum := 0
	for i := p1.X; i <= p2.X; i++ {
		for j := p1.Y; j <= p2.Y; j++ {
			sum += field[i][j]
		}
	}
	return sum
}

func (f *Field) Fill(str []string) {

	for i, s := range str {
		j := 0
		for _, ss := range s {
			if n, err := strconv.Atoi(string(ss)); err == nil {
				f[i][j] = n
				j++
			}
		}
	}
}

func ReadFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadAll(f), nil
}

func ReadAll(file *os.File) []string {
	result := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		os.Exit(1)
	}

	return result
}

func (f Field) Print() {
	// top horizontal line
	fmt.Print(" ")
	for k := 0; k < 23; k++ {
		fmt.Print("-")
	}
	fmt.Println()
	// print sudoku cells
	for i := 0; i < 9; i++ {
		fmt.Print("| ")
		for j := 0; j < 9; j++ {
			fmt.Printf("%-2d", f[i][j])
			if (j+1)%3 == 0 {
				fmt.Print("| ")
			}
		}
		fmt.Println()
		// print horizontal lines
		if (i+1)%3 == 0 {
			fmt.Print(" ")
			for k := 0; k < 23; k++ {
				fmt.Print("-")
			}
			fmt.Println()
		}
	}
}
