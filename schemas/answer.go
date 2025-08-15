package schemas

type Answer struct {
	data [][]int
}

func (a *Answer) InitData(horizontalLength, verticalLength int) {
	a.data = make([][]int, verticalLength)
	for i := range a.data {
		a.data[i] = make([]int, horizontalLength)
	}
}
func (a Answer) GetData() [][]int {
	return a.data
}

func (a Answer) CopyLine(isHorizontal bool, idx int, line *[]int) {
	*line = make([]int, a.GetLength(isHorizontal))
	if isHorizontal {
		copy(*line, a.data[idx])
	} else {
		for i, horizontalLine := range a.data {
			(*line)[i] = horizontalLine[idx]
		}
	}
}
func (a *Answer) SaveLine(isHorizontal bool, idx int, line []int) {
	if len(line) != a.GetLength(isHorizontal) {
		return
	}
	if isHorizontal {
		copy(a.data[idx], line)
	} else {
		for i, v := range line {
			a.data[i][idx] = v
		}
	}
}
func (a Answer) GetLength(isHorizontal bool) int {
	if isHorizontal {
		return len(a.data[0])
	} else {
		return len(a.data)
	}
}

func (a Answer) IsSolved() bool {
	for _, line := range a.data {
		for _, cell := range line {
			if cell == 0 {
				return false
			}
		}
	}
	return true
}
