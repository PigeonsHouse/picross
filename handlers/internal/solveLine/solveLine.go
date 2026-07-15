package solveLine

import (
	"fmt"
	itemrangelist "picross/handlers/internal/solveLine/internal/itemRangeList"
	quizpattern "picross/handlers/internal/solveLine/internal/quizPattern"
	splitanswerline "picross/handlers/internal/solveLine/internal/splitAnswerLine"
	"picross/logger"
	"picross/schemas"
	"slices"
)

type AnswerLine []schemas.CellType

// quizLine: クイズの行 (例: [3, 1, 2])
// answerLine: 解答の行 (例: [_, x, ◼, _, _, x, _])
func (a AnswerLine) SolveLine(quizLine []int) bool {
	if len(quizLine) == 0 {
		panic("quizLineの長さが0")
	}
	// 解き終わっているなら何もしない
	if a.isSolved() {
		return false
	}

	// 現状の解答欄の確定要素を見て、quizLineのどの数字がどの範囲に入りうるか計算する
	itemRangeList := a.getQuizItemRangeList(quizLine)

	// 範囲的にどの数字も届かないマスはUnfilledを付けておく
	isUnfilledUnreachableAreaChanged := a.unfilledUnreachableArea(itemRangeList)

	// 範囲が数字の半分未満であればFilledが確定するので塗る
	// 数字と範囲の長さが一致していれば左右にUnfilledもつける
	isCenterOverlapChanged := a.fillCenterOverlap(itemRangeList)

	// Filledが必要数揃っており解き終わっていたら、残りのセルをUnfilledで埋めておく
	isAutoUnfilledChanged := a.autoUnfilled(quizLine)

	return isCenterOverlapChanged || isAutoUnfilledChanged || isUnfilledUnreachableAreaChanged
}

func (a AnswerLine) isSolved() bool {
	return !slices.Contains(a, schemas.Unsettled)
}

func (a AnswerLine) getQuizItemRangeList(quizLine []int) itemrangelist.ItemRangeList {
	// AnswerLineをUnfilled(x)で分割して、Unsettled(_)とFilled(◼)の塊に分ける
	// Unfilledが連続しても1つ分のUnfilledとみなして分割する
	// [_,_,◼,_,_,x,_,_,x,x,_,_,x] => [[_,_,◼,_,_],[_,_],[_,_]]
	splittedAnswerLine := splitanswerline.SplitAnswerLine(a)
	logger.DebugLog(fmt.Sprintf("xのない部分を分割\n%+v", splittedAnswerLine))

	// splittedAnswerLinesの各要素に、quizLineの要素をどう割り当てられるか、パターンを全探索する
	quizItemPatterns := quizpattern.GenerateQuizPatterns(splittedAnswerLine, quizLine)
	logger.DebugLog(fmt.Sprintf("分割したエリアごとにこの行の問題の数値を割り振る全パターン\n%+v", quizItemPatterns))

	// QuizLineの各数値がAnswerLineのどの範囲に入りうるかの全パターンを計算する
	// そして、start, endを全パターンで比較して、最も広い範囲をitemRangeとする
	itemRangeList := itemrangelist.CalculateItemRangeListPatterns(quizItemPatterns, len(quizLine), splittedAnswerLine)
	logger.DebugLog(fmt.Sprintf("この行の問題の各数値がどの範囲に入りうるか\n%+v", itemRangeList))

	return itemRangeList
}

func (a AnswerLine) fillCenterOverlap(irl itemrangelist.ItemRangeList) bool {
	isChanged := false
	for _, itemRange := range irl {
		// itemRangeの長さが、itemの長さの2倍未満の場合、itemRangeの中央部分は必ず黒になる
		// 例) 2..6の5マスの間に、3を塗る際は、index=4の真ん中1マスだけを塗る
		midStart, midEnd := itemRange.MiddleFilledPosition()
		if midStart <= midEnd {
			for k := midStart; k <= midEnd; k++ {
				if a[k] == schemas.Unsettled {
					a[k] = schemas.Filled
					isChanged = true
				}
				if a[k] == schemas.Unfilled {
					panic("ロジックミス")
				}
			}
		}
		// 範囲の長さと数値が一致している場合は左右にUnfilledを塗る
		if itemRange.IsFit() {
			if midStart > 0 {
				if a[midStart-1] != schemas.Unfilled {
					a[midStart-1] = schemas.Unfilled
					isChanged = true
				}
			}
			if midEnd < len(a)-1 {
				if a[midEnd+1] != schemas.Unfilled {
					a[midEnd+1] = schemas.Unfilled
					isChanged = true
				}
			}
		}
	}
	return isChanged
}

func (a AnswerLine) unfilledUnreachableArea(irl itemrangelist.ItemRangeList) bool {
	isChanged := false
	for unreachableIndex := range irl.UnreachableAreaIndexes(len(a)) {
		if a[unreachableIndex] == schemas.Unsettled {
			a[unreachableIndex] = schemas.Unfilled
			isChanged = true
		}
	}

	if isChanged {
		logger.DebugLog("unreachable debug: isChanged")
	}

	return isChanged
}

func (a AnswerLine) autoUnfilled(quizLine []int) bool {
	// TODO 端から黒確定していたらUnfilledを塗っていく処理も兼ねさせる
	isChanged := false
	quizFilledCount := 0
	answerFilledCount := 0
	for _, quizItem := range quizLine {
		quizFilledCount += quizItem
	}
	for _, cell := range a {
		if cell == schemas.Filled {
			answerFilledCount += 1
		}
	}
	if quizFilledCount == answerFilledCount {
		for i := range a {
			if a[i] == schemas.Unsettled {
				a[i] = schemas.Unfilled
				isChanged = true
			}
		}
	}
	return isChanged
}
