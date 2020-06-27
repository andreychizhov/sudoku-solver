package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Field [9][9]int

type CVector map[int][]Cell

type Cell struct {
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

	data, err := ReadFile("samples/hard/2.txt")
	if err != nil {
		fmt.Print("Error reading from file...")
		os.Exit(1)
	}
	field.Fill(data)

	if err := field.Validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Initial sudoku field")
	field.Print()
	hints := GetHints(field)
	fmt.Printf("Field is correct, %d hints found\n\n", len(hints))

	field.Solve(1)
}

func (f *Field) Solve(step int) {
	if ok, _, _ := f.Check(); ok {
		fmt.Printf("Solution found in %d steps!\n", step-1)
		return
	}

	fmt.Printf("Step %d\n", step)

	areas := GetSquares()
	var hasSolutions = false
	var vectors map[Area]CVector

	for _, a := range areas {
		//fmt.Println(a)
		v := BuildVector(*f, a)
		vectors[a] = v
		s := f.FillArea(v)
		hasSolutions = hasSolutions || s
		PrintVector(v)
	}
	fmt.Println()

	if !hasSolutions {
		// Plan B: trying to resolve ambiguity

		fmt.Println("Sorry, but this sudoku cannot be solved :(")
		return
	}

	f.Print()

	f.Solve(step + 1)
}

func GetMinimalSolution(vrs map[Area]CVector) (Area, CVector) {
	for a, v := range vrs {
		vlen := 0
		for _, c := range v {
			if vlen == 0 || len(c) < vlen {
				vlen = len(c)
			}
		}

	}
}

func (f *Field) FillArea(v CVector) bool {
	var s = false
	for i := 1; i <= 9; i++ {
		if arr, ok := v[i]; ok && len(arr) == 1 {
			x, y := arr[0].X, arr[0].Y
			f[x][y] = i
			s = s || true
		}
	}
	return s
}

func PrintVector(vector CVector) {
	for i := 1; i <= 9; i++ {
		if arr, ok := vector[i]; ok && len(arr) == 1 {
			fmt.Printf("%d: %v; ", i, arr)
		}
	}
}

func BuildVector(f Field, a Area) CVector {
	cands := GetCandidates(f, a)
	vr := make(map[int][]Cell)

	for _, c := range cands {
		for i := a.Tl.X; i <= a.Br.X; i++ {
			for j := a.Tl.Y; j <= a.Br.Y; j++ {
				if f.CanPutIntoCell(i, j, c, false) {
					ar, ok := vr[c]
					if ok {
						vr[c] = append(ar, Cell{i, j})
					} else {
						vr[c] = []Cell{{i, j}}
					}
				}
			}
		}
	}

	return vr
}

func GetCandidates(f Field, a Area) []int {
	all := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	cands := make([]int, 0, 9)

	hints := GetHintsForArea(f, a)

	for _, n := range all {
		if hints[n] {
			continue
		}
		cands = append(cands, n)
	}
	return cands
}

// Hints for area
func GetHintsForArea(f Field, a Area) map[int]bool {
	hints := make(map[int]bool)

	for i := a.Tl.X; i <= a.Br.X; i++ {
		for j := a.Tl.Y; j <= a.Br.Y; j++ {
			if f[i][j] > 0 {
				hints[f[i][j]] = true
			}
		}
	}
	return hints
}

// Hints for field
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

func (f Field) Validate() error {
	hints := GetHints(f)

	if len(hints) < 17 {
		return fmt.Errorf("At least %d hints should be placed on sudoku field.", 17)
	}

	for k, v := range hints {
		if !f.CanPutIntoCell(k.X, k.Y, v, true) {
			return fmt.Errorf("Sudoku field has incorrect initial state, cell %v number %d", k, v)
		}
	}

	return nil
}

func (f Field) CanPutIntoCell(x, y, n int, allowNonEmpty bool) bool {
	if allowNonEmpty {
		f[x][y] = 0
	}

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
	areas := make([]Area, 0, 9)
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
