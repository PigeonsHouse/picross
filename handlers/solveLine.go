package handlers

import (
	"fmt"
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

	if slices.Equal(quizLine, []int{0}) {
		for i := range a {
			a[i] = schemas.Unfilled
		}
		return true
	}

	// AnswerLineをUnfilled(x)で分割して、Unsettled(_)とFilled(◼)の塊に分ける
	// Unfilledが連続しても1つ分のUnfilledとみなして分割する
	// [_,_,◼,_,_,x,_,_,x,x,_,_,x] => [[_,_,◼,_,_],[_,_],[_,_]]
	splittedAnswerLine := a.splitAnswerLine()
	fmt.Println("splitAnswerLine", splittedAnswerLine)

	// splittedAnswerLinesの各要素に、quizLineの要素をどう割り当てられるか、パターンを全探索する
	quizItemPatterns := splittedAnswerLine.generateQuizPatterns(quizLine)
	fmt.Println("generateQuizPatterns", quizItemPatterns)

	// QuizLineの各数値がAnswerLineのどの範囲に入りうるかの全パターンを計算する
	// そして、start, endを全パターンで比較して、最も広い範囲をitemRangeとする
	maxItemRangeList := quizItemPatterns.calculateItemRangeListPatterns(len(quizLine), splittedAnswerLine)
	fmt.Println("calculateItemRangeListPatterns", maxItemRangeList)

	// maxItemRangeの長さが、itemの長さの2倍未満の場合、itemRangeの中央部分は必ず黒になるので、answerLineを更新する
	isCenterOverlapChanged := maxItemRangeList.fillCenterOverlap(a)

	return isCenterOverlapChanged
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

// QuizLineの数字一つ分のデータ
type QuizLineItem struct {
	// QuizLine内のindex
	index int
	// QuizLine内の数値
	value int
}

// SplittedAnswerLineにQuizLineのItemをどう置くかというパターンの一部分
// Lengthは、最大でQuizLineの長さ、最小は0
type QuizItemAllocationInPart []QuizLineItem

// SplittedAnswerLineにQuizLineのItemをどう置くかというパターン
// 常にSplittedAnswerLineの長さと一致する
type QuizItemAllocationPattern []QuizItemAllocationInPart

// 上のパターンのリスト
type QuizItemAllocationPatterns []QuizItemAllocationPattern

func (sal SplittedAnswerLine) generateQuizPatterns(quizLine []int) QuizItemAllocationPatterns {
	// 例: quizLine=[3,1,2], splittedAnswerLines=[[_,_,_,_,_,_,_,_]] の場合、以下のパターンのみとなる
	// - [3,1,2] -> [[3,1,2]]
	// 例: quizLine=[3,1,2], splittedAnswerLines=[[_,_,_,_,_],[_,_],[_,_]] の場合、以下のパターンが考えられる
	// - [3,1,2] -> [[3], [1], [2]]
	// - [3,1,2] -> [[3,1], [2], []]
	// この場合、以下のような値を返す
	// QuizItemPatterns {
	//   QuizItemAllocationPattern {
	//     QuizItemAllocationInPart {
	//       QuizLineItem{index:0, value:3},
	//       QuizLineItem{index:1, value:1},
	//     },
	//     QuizItemAllocationInPart {
	//       QuizLineItem{index:2, value:2},
	//     },
	//     {},
	//   },
	//   QuizItemAllocationPattern {
	//     QuizItemAllocationInPart {
	//       QuizLineItem{index:0, value:3},
	//     },
	//     QuizItemAllocationInPart {
	//       QuizLineItem{index:1, value:1},
	//     },
	//     QuizItemAllocationInPart {
	//       QuizLineItem{index:2, value:2},
	//     },
	//   },
	// }
	// 上の場合は、3や1を1つずついれる場合は3、1それぞれの長さだけで入るが、2つ以上まとめて入れる場合は間のスペースの分も考慮する必要がある
	// 例えば、[3,1]を[_,_,_,_,_]に入れる場合、[3,1]の間に1つスペースが必要なので、3+1+1=5となり、ちょうど入る

	// TODO: まだ仮実装
	quizItemAllocationInPart := make(QuizItemAllocationInPart, len(quizLine))
	for i, item := range quizLine {
		quizItemAllocationInPart[i] = QuizLineItem{
			index: i,
			value: item,
		}
	}
	return QuizItemAllocationPatterns{QuizItemAllocationPattern{quizItemAllocationInPart}}
}

// QuizLineの数値1個分がAnswerLineのうちのどの範囲に入りうるか
type ItemRange struct {
	// itemの左端が入りうる開始位置(全域の場合0が入る)
	start int
	// itemの右端が入りうる終了位置(全域の場合len(answerLine)-1が入る)
	end  int
	item QuizLineItem
}

func (ir ItemRange) Length() int {
	return ir.end - ir.start + 1
}

// ItemRangeのQuizLine一列分
// 長さはQuizLineと一致する
type ItemRangeList []ItemRange

// ItemRangeListの全パターン分
// QuizItemAllocationPatternsと長さは一致する
// 最終的に全パターンを畳み込んで、一番広い範囲のItemRangeListに変換される
type ItemRangeListPatterns []ItemRangeList

// QuizLineの各数値がAnswerLineのどの範囲に入りうるかの全パターンを計算する
// そしてその各パターンから、各QuizLineItemの取りうる最大の範囲を計算する
func (qiap QuizItemAllocationPatterns) calculateItemRangeListPatterns(quizLineLength int, splittedAnswerLine SplittedAnswerLine) ItemRangeList {
	// 各QuizLineItemの割り振りパターンでのQuizLineItemの各入りうる範囲
	itemRangeListPatterns := make(ItemRangeListPatterns, len(qiap))

	// QuizLineItemの割り振りパターンでforを回す
	for h, itemAllocationPattern := range qiap {
		itemRangeList := make(ItemRangeList, 0, quizLineLength)
		offset := 0
		// その中で、splitされた一箇所の、QuizLineItemとAnswerの値を取り出す
		for i := range splittedAnswerLine {
			// TODO: answerにFilledがあって取りうる範囲が変わるケースにも対応する

			// QuizLineの一部分。構造体なので、全体のQuizLineの何番目かのindexも取れる
			// [3,1,2]
			partQuizLineItem := itemAllocationPattern[i]
			// AnswerLineの一部分
			// [_,_,◼,_,_,_,_,_,_]
			answer := splittedAnswerLine[i]

			// QuizLineの数値ひとつに注目
			for j, quizItem := range partQuizLineItem {
				// この周回のquizItemを除いた、左側の数値の合計+余白マス数分startをずらし、右側の数値の合計+余白マス数分endをずらして、範囲を計算する
				start := 0
				for k := range j {
					start += partQuizLineItem[k].value + 1
				}
				end := len(answer) - 1
				for k := j + 1; k < len(partQuizLineItem); k++ {
					end -= partQuizLineItem[k].value + 1
				}

				itemRangeList = append(itemRangeList, ItemRange{
					start: offset + start,
					end:   offset + end,
					item:  quizItem,
				})
			}

			// 次の周回で、これまでのsplitのマス数を把握するためoffsetに加算
			offset += len(answer)
		}
		// 1パターン分できたのでpush
		itemRangeListPatterns[h] = itemRangeList
	}
	fmt.Println("itemRangeListPatterns", itemRangeListPatterns)

	maxItemRangeList := make(ItemRangeList, quizLineLength)
	// 各パターンを回し、各QuizLineItemの入りうる範囲の一番広い範囲を計算する
	for h, itemRangeList := range itemRangeListPatterns {
		for i, itemRange := range itemRangeList {
			if h == 0 {
				maxItemRangeList[i] = itemRange
			} else {
				maxItemRangeList[i].start = min(itemRange.start, maxItemRangeList[i].start)
				maxItemRangeList[i].end = max(itemRange.end, maxItemRangeList[i].end)
			}
		}
	}
	return maxItemRangeList
}

func (irl ItemRangeList) fillCenterOverlap(answerLine AnswerLine) bool {
	isChanged := false
	for _, itemRange := range irl {
		// itemRangeの長さが、itemの長さの2倍未満の場合、itemRangeの中央部分は必ず黒になる
		itemRangeLength := itemRange.Length()

		// 例) 2..6の5マスの間に、3を塗る際は、index=4の真ん中1マスだけを塗る
		unsettleLength := itemRangeLength - itemRange.item.value
		midStart := itemRange.start + unsettleLength
		midEnd := itemRange.end - unsettleLength
		if midStart <= midEnd {
			for k := midStart; k <= midEnd; k++ {
				if answerLine[k] == schemas.Unsettled {
					answerLine[k] = schemas.Filled
					isChanged = true
				}
				if answerLine[k] == schemas.Unfilled {
					panic("ロジックミス")
				}
			}
		}
		// 範囲の長さと数値が一致している場合は左右にUnfilledを塗る
		if unsettleLength == 0 {
			if midStart > 0 {
				answerLine[midStart-1] = schemas.Unfilled
			}
			if midEnd < len(answerLine)-1 {
				answerLine[midEnd+1] = schemas.Unfilled
			}
		}
	}
	return isChanged
}
