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

	// AnswerLineをUnfilled(x)で分割して、Unsettled(_)とFilled(◼)の塊に分ける
	// Unfilledが連続しても1つ分のUnfilledとみなして分割する
	// [_,_,◼,_,_,x,_,_,x,x,_,_,x] => [[_,_,◼,_,_],[_,_],[_,_]]
	splittedAnswerLine := a.splitAnswerLine()
	fmt.Printf("xのない部分を分割\n　　%+v\n", splittedAnswerLine)

	// splittedAnswerLinesの各要素に、quizLineの要素をどう割り当てられるか、パターンを全探索する
	quizItemPatterns := splittedAnswerLine.generateQuizPatterns(quizLine)
	fmt.Printf("分割したエリアごとにこの行の問題の数値を割り振る全パターン\n　　%+v\n", quizItemPatterns)

	// QuizLineの各数値がAnswerLineのどの範囲に入りうるかの全パターンを計算する
	// そして、start, endを全パターンで比較して、最も広い範囲をitemRangeとする
	maxItemRangeList := quizItemPatterns.calculateItemRangeListPatterns(len(quizLine), splittedAnswerLine)
	fmt.Printf("この行の問題の各数値がどの範囲に入りうるか\n　　%+v\n", maxItemRangeList)

	// maxItemRangeの長さが、itemの長さの2倍未満の場合、itemRangeの中央部分は必ず黒になるので、answerLineを更新する
	isCenterOverlapChanged := a.fillCenterOverlap(maxItemRangeList)

	isAutoUnfilledChanged := a.autoUnfilled(quizLine)

	isUnfilledUnreachableAreaChanged := a.unfilledUnreachableArea(maxItemRangeList)

	return isCenterOverlapChanged || isAutoUnfilledChanged || isUnfilledUnreachableAreaChanged
}

func (a AnswerLine) isSolved() bool {
	return !slices.Contains(a, schemas.Unsettled)
}

// AnswerLineをUnfilled(x)で分割した一部分
// Unsettled(_)とFilled(◼)だけが含まれている
// Unfilledを無視する関係で、startとendを保持している
type SplittedAnswerLinePart struct {
	start int
	end   int
	cells []schemas.CellType
}

func (salp SplittedAnswerLinePart) String() string {
	return fmt.Sprintf("{start:%d,end:%d,cells:%+v}", salp.start, salp.end, salp.cells)
}

// AnswerLineをUnfilled(x)で分割したデータ
// 例えばUnfilledのないAnswerLineはlength=1になる
type SplittedAnswerLine []SplittedAnswerLinePart

// AnswerLineをUnfilled(x)で分割する
func (a AnswerLine) splitAnswerLine() SplittedAnswerLine {
	result := SplittedAnswerLine{}
	currentPart := SplittedAnswerLinePart{}

	for i, cell := range a {
		if cell == schemas.Unfilled {
			if len(currentPart.cells) > 0 {
				result = append(result, currentPart)
				currentPart = SplittedAnswerLinePart{}
			}
			currentPart.start = i + 1
		} else {
			currentPart.cells = append(currentPart.cells, cell)
			currentPart.end = i
		}
	}
	if len(currentPart.cells) > 0 {
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

func convertQuizLineItemList(quizLine []int) []QuizLineItem {
	quizLineItems := make([]QuizLineItem, len(quizLine))
	for i, num := range quizLine {
		quizLineItems[i] = QuizLineItem{
			index: i,
			value: num,
		}
	}
	return quizLineItems
}

func (sal SplittedAnswerLine) generateQuizPatterns(quizLine []int) QuizItemAllocationPatterns {
	// 例: quizLine=[3,1,2], sal=[[_,_,_,_,_,_,_,_]] の場合、以下のパターンのみとなる
	// - [3,1,2] -> [[3,1,2]]
	// 例: quizLine=[3,1,2], sal=[[_,_,_,_,_],[_,_],[_,_]] の場合、以下のパターンが考えられる
	// - [3,1,2] -> [[3], [1], [2]]
	// - [3,1,2] -> [[3,1], [2], []]
	// - [3,1,2] -> [[3,1], [], [2]]
	// この場合、以下のような値を返す
	// QuizItemPatterns {
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
	//       QuizLineItem{index:1, value:1},
	//     },
	//     {},
	//     QuizItemAllocationInPart {
	//       QuizLineItem{index:2, value:2},
	//     },
	//   },
	// }
	// 上の場合は、3や1を1つずついれる場合は3、1それぞれの長さだけで入るが、2つ以上まとめて入れる場合は間のスペースの分も考慮する必要がある
	// 例えば、[3,1]を[_,_,_,_,_]に入れる場合、[3,1]の間に1つスペースが必要なので、3+1+1=5となり、ちょうど入る

	workingQuizLineItems := convertQuizLineItemList(quizLine)

    isContainable := func (quizLine []QuizLineItem, cells []schemas.CellType) bool {
		sum := -1
		for _, item := range quizLine {
			sum += item.value+1
		}
		if sum == len(cells) {
			idx := -1
			for _, item := range quizLine {
				idx += item.value+1
				if idx == len(cells) {
					break
				}
				if cells[idx] == schemas.Filled {
					return false
				}
			}
			return true
		} else {
			return sum<len(cells)
		}
	}

	var recFunc func(innerQuizLine []QuizLineItem, innerSal SplittedAnswerLine) QuizItemAllocationPatterns
	recFunc = func(innerQuizLine []QuizLineItem, innerSal SplittedAnswerLine) QuizItemAllocationPatterns {
		// fmt.Printf("start recFunc\n　quiz: %+v target answer: %+v\n", innerQuizLine, innerSal[0])
		var wholePatterns QuizItemAllocationPatterns
		if len(innerQuizLine) == 0 {
			// innerSalの長さ分partを含める必要がある
			wholePatterns = append(wholePatterns,
				make(QuizItemAllocationPattern, len(innerSal)),
			)
		} else if len(innerSal) == 1 {
			if isContainable(innerQuizLine, innerSal[0].cells) {
				wholePatterns = append(wholePatterns,
					QuizItemAllocationPattern{
						QuizItemAllocationInPart(innerQuizLine),
					},
				)
			} else {
				// 入り切らないケースはnilとしてしまう
				// fmt.Printf(
				// 	"finish recFunc\n　quiz: %+v target answer: %+v result: nil\n",
				// 	innerQuizLine, innerSal[0],
				// )
				return nil
			}
		} else {
			// fmt.Println("current part: []")
			patterns := recFunc(innerQuizLine, innerSal[1:])
			if patterns != nil {
				for i := range patterns {
					patterns[i] = append(
						QuizItemAllocationPattern{QuizItemAllocationInPart{}},
						patterns[i]...,
					)
				}
				// fmt.Printf("append pattern: %+v\n", patterns)
				wholePatterns = append(wholePatterns, patterns...)
			}

			for i := range innerQuizLine {
				target := innerQuizLine[:i+1]
				if !isContainable(target, innerSal[0].cells) {
					break
				}
				// fmt.Println("current part:", innerQuizLine[:i+1])
				patterns = recFunc(innerQuizLine[i+1:], innerSal[1:])
				if patterns != nil {
					for j := range patterns {
						patterns[j] = append(
							QuizItemAllocationPattern{
								QuizItemAllocationInPart(target),
							},
							patterns[j]...,
						)
					}
					// fmt.Printf("append pattern: %+v\n", patterns)
					wholePatterns = append(wholePatterns, patterns...)
				}
			}
		}
		// fmt.Printf("finish recFunc\n　quiz: %+v target answer: %+v result: %+v\n", innerQuizLine, innerSal[0], wholePatterns)
		return wholePatterns
	}
	return recFunc(workingQuizLineItems, sal)
}

// QuizLineの数値1個分がAnswerLineのうちのどの範囲に入りうるか
type ItemRange struct {
	// itemの左端が入りうる開始位置(全域の場合0が入る)
	start int
	// itemの右端が入りうる終了位置(全域の場合len(answerLine)-1が入る)
	end int

	filledStart *int
	filledEnd   *int

	item QuizLineItem
}

func (ir ItemRange) String() string {
	if ir.filledStart != nil {
		return fmt.Sprintf(
			"{start:%d end:%d filledStart:%d filledEnd:%d item:%+v}",
			ir.start, ir.end, *ir.filledStart, *ir.filledEnd, ir.item,
		)
	} else {
		return fmt.Sprintf("{start:%d end:%d item:%+v}", ir.start, ir.end, ir.item)
	}
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
		// その中で、splitされた一箇所の、QuizLineItemとAnswerの値を取り出す
		for i := range splittedAnswerLine {
			// QuizLineの一部分。構造体なので、全体のQuizLineの何番目かのindexも取れる
			// [3,1,2]
			partQuizLineItem := itemAllocationPattern[i]
			// AnswerLineの一部分
			// [_,_,_,◼,_,_,_,_,_]
			answer := splittedAnswerLine[i]

			// fmt.Printf("=====\n　確認中の問題の一部: %+v\n　その数値が入る解答欄の一部: %+v\n", partQuizLineItem, answer)

			// QuizLineの数値ひとつに注目
			for j, quizItem := range partQuizLineItem {
				// この周回のquizItemを除いた、左側の数値の合計+余白マス数分startをずらし、右側の数値の合計+余白マス数分endをずらして、範囲を計算する
				start := 0
				for k := range j {
					start += partQuizLineItem[k].value + 1
				}
				end := len(answer.cells) - 1
				for k := j + 1; k < len(partQuizLineItem); k++ {
					end -= partQuizLineItem[k].value + 1
				}
				// start,endはグローバルな一次元座標なので、オフセットとなるanswer.startを加算している
				ir := ItemRange{
					start: answer.start + start,
					end:   answer.start + end,
					item:  quizItem,
				}
				// 一つ前のitemRangeのfilledEndよりもstartが左にある場合、手放す必要があるのでずらす
				if j != 0 && itemRangeList[len(itemRangeList)-1].filledEnd != nil {
					ir.start = max(ir.start, *itemRangeList[len(itemRangeList)-1].filledEnd+1)
				}
				// fmt.Printf("answer確認前の範囲: %d ~ %d\n", ir.start, ir.end)

				// startとendの中にFilledがある場合はそれを中心に範囲を狭める
				// FIXME: このロジックがおかしい
				filledBlockStartIndex := -1
				filledBlockEndIndex := -1
				// fmt.Printf("answer切り出し: %+v\n", answer.cells[ir.start-answer.start:ir.end-answer.start+1])
				for k, cell := range answer.cells[ir.start-answer.start : ir.end-answer.start+1] {
					if cell == schemas.Filled {
						if filledBlockStartIndex == -1 {
							filledBlockStartIndex = k + ir.start
						}
						filledBlockEndIndex = k + ir.start
					} else if filledBlockEndIndex != -1 {
						break
					}
				}
				// fmt.Println(ir.start, ir.end, filledBlockStartIndex, filledBlockEndIndex, quizItem.value)
				// 違うfilledの塊を掴んでいたら手放す
				if quizItem.value < filledBlockEndIndex-filledBlockStartIndex+1 {
					filledBlockStartIndex = -1
					filledBlockEndIndex = -1
				}
				// filledがあればrangeを調整
				if filledBlockStartIndex != -1 {
					ir.start = max(ir.start, filledBlockEndIndex-(quizItem.value-1))
					ir.end = min(ir.end, filledBlockStartIndex+(quizItem.value-1))
					ir.filledStart = &filledBlockStartIndex
					ir.filledEnd = &filledBlockEndIndex
				}
				// fmt.Printf("answer確認後の範囲: %d ~ %d\n", ir.start, ir.end)
				// fmt.Printf("最終的なItemRange: %+v\n", ir)

				itemRangeList = append(itemRangeList, ir)
			}

			// fmt.Printf("各数値が入る範囲: %+v\n", itemRangeList)
		}
		// 1パターン分できたのでpush
		itemRangeListPatterns[h] = itemRangeList
	}
	// fmt.Printf("各パターンごとの数値の取りうる範囲\n　　%+v\n", itemRangeListPatterns)

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

func (a AnswerLine) fillCenterOverlap(irl ItemRangeList) bool {
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
		if unsettleLength == 0 {
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

func (a AnswerLine) unfilledUnreachableArea(irl ItemRangeList) bool {
	isChanged := false
	checkedIndex := -1
	for h, itemRange := range irl {
		if checkedIndex+1 < itemRange.start {
			fmt.Println("unreachable debug:", h, a[checkedIndex+1:itemRange.start])
			for i := range a[checkedIndex+1 : itemRange.start] {
				if a[checkedIndex+1+i] == schemas.Unsettled {
					a[checkedIndex+1+i] = schemas.Unfilled
					isChanged = true
				}
			}
		}
		checkedIndex = itemRange.end
	}
	if checkedIndex+1 < len(a) {
		fmt.Println("unreachable debug: last", a[checkedIndex+1:])
		for i := range a[checkedIndex+1:] {
			if a[checkedIndex+1+i] == schemas.Unsettled {
				a[checkedIndex+1+i] = schemas.Unfilled
				isChanged = true
			}
		}
	}

	if isChanged {
		fmt.Println("unreachable debug: isChanged")
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
