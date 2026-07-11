package handlers

import (
	"fmt"
	"picross/schemas"
	"slices"
)

type AnswerLine []schemas.CellType

// quizLine: クイズの行 (例: [3, 1, 2])
// answerLine: 解答の行 (例: [_, x, ◼, _, _, x, _])
func (a AnswerLine) SolveLine(quizLine []int) (isChanged bool) {
	if len(quizLine) == 0 {
		panic("quizLineの長さが0")
	}
	// 解き終わっているなら何もしない
	if a.isSolved() {
		return
	}

	// answerLineを-1でsplitして、0と1の塊に分ける(-1の連続は1つにまとめる)
	splittedAnswerLine := a.splitAnswerLine()

	// splittedAnswerLinesの各要素に、quizLineの要素をどう割り当てられるか、パターンを全探索する
	quizPatterns := splittedAnswerLine.generateQuizPatterns(quizLine)
	fmt.Println("splittedAnswerLines", splittedAnswerLine, "patterns:", quizPatterns)

	itemRangeListPatterns := make([][]ItemRange, len(quizPatterns), len(splittedAnswerLine[0]))
	for h, splitedQuiz := range quizPatterns {
		for i := range splittedAnswerLine {
			partQuiz := splitedQuiz[i]
			answer := splittedAnswerLine[i]
			itemRangeListPatterns[h] = make([]ItemRange, len(partQuiz))

			for j, quizItem := range partQuiz {
				// quizItem以外の、左のvalueの合計+item数-1、右のvalueの合計+item数-1を、左右から引いて、quizItemの入りうる位置を特定する
				leftMin := 0
				for k := 0; k < j; k++ {
					leftMin += partQuiz[k].value + 1
				}
				rightMin := 0
				for k := j + 1; k < len(partQuiz); k++ {
					rightMin += partQuiz[k].value + 1
				}

				// TODO: answerの値をintからQuizItemにして、indexとvalueが同じQuizItemの場合はその黒から届く範囲がitemRangeとする
				itemRangeListPatterns[h][j] = ItemRange{
					start: leftMin,
					end:   len(answer) - 1 - rightMin,
					item:  quizItem,
				}
			}
		}
	}
	// start, endを全パターンで比較して、最も広い範囲をitemRangeとする
	maxItemRangeList := make([]ItemRange, len(quizLine))
	for _, itemRangeList := range itemRangeListPatterns {
		for i, itemRange := range itemRangeList {
			if itemRange.start < maxItemRangeList[i].start {
				maxItemRangeList[i].start = itemRange.start
			}
			if itemRange.end > maxItemRangeList[i].end {
				maxItemRangeList[i].end = itemRange.end
			}
		}
	}
	// maxItemRangeの長さが、itemの長さの2倍未満の場合、itemRangeの中央部分は必ず黒になるので、answerLineを更新する
	for _, maxItemRange := range maxItemRangeList {
		itemRangeLength := maxItemRange.Length()
		if itemRangeLength < maxItemRange.item.value*2 {
			// itemRangeの長さが、itemの長さの2倍未満の場合、itemRangeの中央部分は必ず黒になる
			midStart := maxItemRange.end - maxItemRange.item.value
			midEnd := maxItemRange.start + maxItemRange.item.value - 1
			for k := midStart; k <= midEnd; k++ {
				if a[k] != schemas.Filled {
					a[k] = schemas.Filled
					isChanged = true
				}
			}
		}
	}

	return
}

func (a AnswerLine) isSolved() bool {
	return !slices.Contains(a, schemas.Unsettled)
}

// AnswerLineをUnfilled(x)で分割した一部分
// Unsettled(_)とFilled(◼)だけが含まれている
type SplittedAnswerLinePart []schemas.CellType

// AnswerLineをUnfilled(x)で分割したデータ
// 例えばUnfilledのないAnswerLineはlength=1になる
type SplittedAnswerLine []SplittedAnswerLinePart

// AnswerLineをUnfilled(x)で分割する
func (a AnswerLine) splitAnswerLine() SplittedAnswerLine {
	result := SplittedAnswerLine{}
	currentPart := SplittedAnswerLinePart{}

	for _, cell := range a {
		if cell == schemas.Unfilled {
			if len(currentPart) > 0 {
				result = append(result, currentPart)
				currentPart = SplittedAnswerLinePart{}
			}
		} else {
			currentPart = append(currentPart, cell)
		}
	}
	if len(currentPart) > 0 {
		result = append(result, currentPart)
	}

	return result
}

type QuizItem struct {
	index int
	value int
}
type ItemRange struct {
	start int // itemの左端が入りうる開始位置(全域の場合0が入る)
	end   int // itemの右端が入りうる終了位置(全域の場合len(answerLine)-1が入る)
	item  QuizItem
}

func (ir ItemRange) Length() int {
	return ir.end - ir.start + 1
}

func (sal SplittedAnswerLine) generateQuizPatterns(quizLine []int) [][][]QuizItem {
	// 例: quizLine=[3,1,2], splittedAnswerLines=[[_,_,_,_,_,_,_,_]] の場合、以下のパターンのみとなる
	// - [3,1,2] -> [[3,1,2]]
	// 例: quizLine=[3,1,2], splittedAnswerLines=[[_,_,_,_,_],[_,_],[_,_]] の場合、以下のパターンが考えられる
	// - [3,1,2] -> [[3], [1], [2]]
	// - [3,1,2] -> [[3,1], [2], []]
	// この場合、以下のような値を返す
	// [][][]QuizItem {
	//   {
	//     {
	//       QuizItem{index:0, value:3},
	//       QuizItem{index:1, value:1},
	//     },
	//     {
	//       QuizItem{index:2, value:2},
	//     },
	//     {},
	//   },
	//   {
	//     {
	//       QuizItem{index:0, value:3},
	//     },
	//     {
	//       QuizItem{index:1, value:1},
	//     },
	//     {
	//       QuizItem{index:2, value:2},
	//     },
	//   },
	// }
	// 上の場合は、3や1を1つずついれる場合は3、1それぞれの長さだけで入るが、2つ以上まとめて入れる場合は間のスペースの分も考慮する必要がある
	// 例えば、[3,1]を[_,_,_,_,_]に入れる場合、[3,1]の間に1つスペースが必要なので、3+1+1=5となり、ちょうど入る

	quizItemList := make([]QuizItem, len(quizLine))
	for idx, value := range quizLine {
		quizItemList[idx] = QuizItem{
			index: idx,
			value: value,
		}
	}
	for {
		tempList := make([][]QuizItem, len(sal))
		tempList[0] = quizItemList

		break
	}

	// TODO: 実装
	return [][][]QuizItem{
		{
			{},
		},
	}
}
