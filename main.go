package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Field [9][9]int

type Point struct { // Cell
	X int
	Y int
}

type Area struct {
	Tl Point
	Br Point
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
	fmt.Print(len(GetHints(field)))

}

func (f *Field) Solve() {
	var subs Point

	hints := GetHints(f)

}

func GetHints(f Field) (map[Point]int) {
	var hints map[Point]int
	hints = make(map[Point]int)
	
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if f[i][j] > 0 {
				hints[Point{i, j}] = f[i][j]
			}
		}
	}
	return hints
}

func (f Field) CanPutIntoCell(x, y, n int) bool {
	var areas = GetAreas(Point(x, y))

	for _, a := range areas {
		if ExistsInArea(f, a, n) {
			return false;
		}
	}
	return true;
}

func ExistsInArea(f Field, a Area, n int) : bool{
	for i := a.Tl.X; i <= a.Br.X; i++ {
		for j := a.Tl.Y; j <= a.Br.Y; j++ {
			if (f[i, j] == n){
				return true;
			}
		}
	}
	return false;
}

func GetAreas(p Point) []Area {
	var areas []Area

	append(areas, Area{Point{p.X, 0}, Point{p.X, 8}})
	append(areas, Area{Point{0, p.Y}, Point{8, p.Y}})

	for i := 0; i <= 6; i += 3 {
		for j := 0; j <= 6; j += 3 {
				if p.X >= i && p.X <= i + 2 && p.Y >= j && p.Y <= j + 2
				{
					append(areas, Area{Point{i, j}, Point{i + 2, j + 2}})
				}
			}
		}
	}
}

func (f *Field) Check() (bool, Point, int) {

	for i := 0; i < 9; i++ {
		if s := Sum(*f, Point{i, 0}, Point{i, 8}); s != 45 {
			return false, Point{i, 8}, s
		}
		if s := Sum(*f, Point{0, i}, Point{8, i}); s != 45 {
			return false, Point{8, i}, s
		}
	}

	for i := 0; i <= 6; i += 3 {
		for j := 0; j <= 6; j += 3 {
			if s := Sum(*f, Point{i, j}, Point{i + 2, j + 2}); s != 45 {
				return false, Point{i, j}, s
			}
		}
	}

	return true, Point{0, 0}, 45
}

func Sum(field Field, p1, p2 Point) int {
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
