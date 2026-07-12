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

	isAutoUnfilledBeforeChanged := a.autoUnfilled(quizLine)

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
	isCenterOverlapChanged := maxItemRangeList.fillCenterOverlap(a)

	// isAutoUnfilledAfterChanged := a.autoUnfilled(quizLine)

	return isCenterOverlapChanged || isAutoUnfilledBeforeChanged
}

func (a AnswerLine) isSolved() bool {
	return !slices.Contains(a, schemas.Unsettled)
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

	quizItemAllocationPattern := make(QuizItemAllocationPattern, len(sal))

	workingQuizLineItems := convertQuizLineItemList(quizLine)
	for i, part := range sal {
		partLength := len(part.cells)

		// 1周目は余白の+1が不要だが、条件分岐をサボるため初期値を-1にして誤魔化している
		fillingQuizItemLengthInPart := -1
		// このパートに含める切り出しのインデックス
		fillingQuizItemIndex := 0
		for j, item := range workingQuizLineItems {
			// はみ出したらj番目は含めない
			if fillingQuizItemLengthInPart+item.value+1 > partLength {
				break
			}
			// 前のquizLineItemとの余白の1マスを明けるため+1をしている
			fillingQuizItemLengthInPart += item.value + 1
			// インデックス切り出しは[l:h]のh-1まで取得するため、indexを1大きくしておく
			fillingQuizItemIndex = j + 1

			// 例えば[_,_,_,◼,_]に[3,1]を割り振る際、ぴったり詰めると間もFilledになってしまう
			// そのため、左寄せで配置した後、その右に続くFilled分ずらす必要がある
			k := 0
			for k < len(part.cells)-1-fillingQuizItemLengthInPart {
				if part.cells[fillingQuizItemLengthInPart+k+1] == schemas.Filled {
					k += 1
				} else {
					break
				}
			}
			fillingQuizItemLengthInPart += k
		}

		// はみ出さなかったitemを今パターンのPartとして追加、切り出した残りは引き続きworkingとして使う
		quizItemAllocationPattern[i] = workingQuizLineItems[:fillingQuizItemIndex]
		workingQuizLineItems = workingQuizLineItems[fillingQuizItemIndex:]
	}

	// TODO: 現状左端に詰めたパターンのみ。他のパターンも計算する

	return QuizItemAllocationPatterns{quizItemAllocationPattern}
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
		// その中で、splitされた一箇所の、QuizLineItemとAnswerの値を取り出す
		for i := range splittedAnswerLine {
			// QuizLineの一部分。構造体なので、全体のQuizLineの何番目かのindexも取れる
			// [3,1,2]
			partQuizLineItem := itemAllocationPattern[i]
			// AnswerLineの一部分
			// [_,_,_,◼,_,_,_,_,_]
			answer := splittedAnswerLine[i]

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

				// startとendの中にFilledがある場合はそれを中心に範囲を狭める
				// 実装メモ
				// [3,1,2]
				// [_,_,_,◼,_,_,_,_,_]
				// 3は ls=0,le=3 k_index=3の段階でfilledBlockが終わり
				//
				// filledBlockEndIndex - itemRange.item.valueとstartの内大きい(右寄り)なindexをstartとする
				// filledBlockStartIndex + itemRange.item.valueとendの内小さい(左寄り)なindexをendとする
				// unsettledが来る前に切り出したanswerが終わるケースにも注意
				// blockが2つ含まれるケースはまた後で
				filledBlockStartIndex := -1
				filledBlockEndIndex := -1
				for k, cell := range answer.cells[start : end+1] {
					if cell == schemas.Filled {
						if filledBlockStartIndex == -1 {
							filledBlockStartIndex = k
						}
						filledBlockEndIndex = k
					}
				}
				// filledがあればrangeを調整
				if filledBlockStartIndex != -1 {
					start = max(start, filledBlockEndIndex-quizItem.value)
					end = max(start, filledBlockStartIndex+quizItem.value)
				}

				// start,endはグローバルな一次元座標なので、オフセットとなるanswer.startを加算している
				itemRangeList = append(itemRangeList, ItemRange{
					start: answer.start + start,
					end:   answer.start + end,
					item:  quizItem,
				})
			}

			fmt.Printf("=====\n　確認中の問題の一部: %+v\n　その数値が入る解答欄の一部: %+v\n　各数値が入る範囲: %+v\n", partQuizLineItem, answer, itemRangeList)
		}
		// 1パターン分できたのでpush
		itemRangeListPatterns[h] = itemRangeList
	}
	fmt.Printf("各パターンごとの数値の取りうる範囲\n　　%+v\n", itemRangeListPatterns)

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
