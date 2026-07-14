package splitanswerline

import (
	"fmt"
	"picross/schemas"
)

// AnswerLineをUnfilled(x)で分割した一部分
// Unsettled(_)とFilled(◼)だけが含まれている
// Unfilledを無視する関係で、startとendを保持している
type SplittedAnswerLinePart struct {
	Start int
	End   int
	Cells []schemas.CellType
}

func (salp SplittedAnswerLinePart) Length() int {
	return len(salp.Cells)
}
func (salp SplittedAnswerLinePart) String() string {
	return fmt.Sprintf("{start:%d,end:%d,cells:%+v}", salp.Start, salp.End, salp.Cells)
}

func (sal SplittedAnswerLinePart) IsContainable(quizLine []int) bool {
	sum := -1
	for _, value := range quizLine {
		sum += value + 1
	}
	if sum == len(sal.Cells) {
		idx := -1
		for _, value := range quizLine {
			idx += value + 1
			if idx == len(sal.Cells) {
				break
			}
			if sal.Cells[idx] == schemas.Filled {
				return false
			}
		}
		return true
	} else {
		return sum < len(sal.Cells)
	}
}

// AnswerLineをUnfilled(x)で分割したデータ
// 例えばUnfilledのないAnswerLineはlength=1になる
type SplittedAnswerLine []SplittedAnswerLinePart

// AnswerLineをUnfilled(x)で分割する
func SplitAnswerLine(answerLine []schemas.CellType) SplittedAnswerLine {
	result := SplittedAnswerLine{}
	currentPart := SplittedAnswerLinePart{}

	for i, cell := range answerLine {
		if cell == schemas.Unfilled {
			if len(currentPart.Cells) > 0 {
				result = append(result, currentPart)
				currentPart = SplittedAnswerLinePart{}
			}
			currentPart.Start = i + 1
		} else {
			currentPart.Cells = append(currentPart.Cells, cell)
			currentPart.End = i
		}
	}
	if len(currentPart.Cells) > 0 {
		result = append(result, currentPart)
	}

	return result
}
