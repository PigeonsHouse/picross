package schemas

import (
	"fmt"
	"strings"
)

type CellType int

const (
	// 未確定(_)
	Unsettled CellType = iota
	// 塗りつぶし(◼)
	Filled
	// 空白(x)
	Unfilled
)

func (c CellType) String() string {
	switch c {
	case Unsettled:
		return "_"
	case Filled:
		return "◼"
	case Unfilled:
		return "x"
	default:
		panic("存在しないCellTypeです")
	}
}

type Answer struct {
	cells [][]CellType
}

type LogOption struct {
	Orientation Orientation
	Index       int
}

func (a Answer) Log(option *LogOption) {
	if option == nil {
		for _, line := range a.cells {
			fmt.Println(line)
		}
	} else {
		spaces := strings.Repeat(" ", a.GetLength(Horizontal)*2+3)
		if option.Orientation == Vertical {
			runes := []rune(spaces)
			runes[(option.Index+1)*2+1] = 'v'
			spaces = string(runes)
		}
		fmt.Printf("%s\n", spaces)
		for i, line := range a.cells {
			cursor := "  "
			if i == option.Index && option.Orientation == Horizontal {
				cursor = "> "
			}
			fmt.Printf("%s%v\n", cursor, line)
		}
	}
}

func (a *Answer) Initialize(horizontalLength, verticalLength int) {
	a.cells = make([][]CellType, verticalLength)
	for i := range a.cells {
		a.cells[i] = make([]CellType, horizontalLength)
	}
}

func (a Answer) Map(callback func(x, y int, cell CellType)) {
	for y, inner := range a.cells {
		for x, data := range inner {
			callback(x, y, data)
		}
	}
}

func (a Answer) ReadLine(orientation Orientation, idx int) []CellType {
	line := make([]CellType, a.GetLength(orientation))
	if orientation == Horizontal {
		copy(line, a.cells[idx])
	} else {
		for i, horizontalLine := range a.cells {
			(line)[i] = horizontalLine[idx]
		}
	}
	return line
}

func (a *Answer) WriteLine(orientation Orientation, idx int, line []CellType) {
	if len(line) != a.GetLength(orientation) {
		return
	}
	if orientation == Horizontal {
		copy(a.cells[idx], line)
	} else {
		for i, cell := range line {
			a.cells[i][idx] = cell
		}
	}
}

func (a Answer) GetLength(orientation Orientation) int {
	if orientation == Horizontal {
		return len(a.cells[0])
	} else {
		return len(a.cells)
	}
}

func (a Answer) IsSolved() bool {
	isSolved := true

	a.Map(func(_, _ int, cell CellType) {
		if cell == Unsettled {
			isSolved = false
		}
	})

	return isSolved
}
