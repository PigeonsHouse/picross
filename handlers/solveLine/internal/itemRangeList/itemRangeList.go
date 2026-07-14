package itemrangelist

import (
	"fmt"
	"iter"
	quizpattern "picross/handlers/solveLine/internal/quizPattern"
	splitanswerline "picross/handlers/solveLine/internal/splitAnswerLine"
	"picross/schemas"
)

// QuizLineの数値1個分がAnswerLineのうちのどの範囲に入りうるか
type ItemRange struct {
	// itemの左端が入りうる開始位置(全域の場合0が入る)
	start int
	// itemの右端が入りうる終了位置(全域の場合len(answerLine)-1が入る)
	end int

	// FIXME: これいらないかも
	filledStart *int
	filledEnd   *int

	item quizpattern.QuizLineItem
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
func (ir ItemRange) MiddleFilledPosition() (int, int) {
	// 例) 2..6の5マスの間に、3を塗る際は、index=4の真ん中1マスだけを塗る
	itemRangeLength := ir.end - ir.start + 1

	unsettleLength := itemRangeLength - ir.item.Value

	midStart := ir.start + unsettleLength
	midEnd := ir.end - unsettleLength
	return midStart, midEnd
}
func (ir ItemRange) IsFit() bool {
	itemRangeLength := ir.end - ir.start + 1
	return itemRangeLength == ir.item.Value
}

// ItemRangeのQuizLine一列分
// 長さはQuizLineと一致する
type ItemRangeList []ItemRange

func (irl ItemRangeList) UnreachableAreaIndexes(length int) iter.Seq[int] {
	return func(yield func(int) bool) {
		checkedIndex := -1
		for _, itemRange := range irl {
			for i := checkedIndex + 1; i < itemRange.start; i++ {
				if !yield(i) {
					return
				}
			}
			checkedIndex = itemRange.end
		}
		for i := checkedIndex + 1; i < length; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// ItemRangeListの全パターン分
// QuizItemAllocationPatternsと長さは一致する
// 最終的に全パターンを畳み込んで、一番広い範囲のItemRangeListに変換される
type ItemRangeListPatterns []ItemRangeList

// QuizLineの各数値がAnswerLineのどの範囲に入りうるかの全パターンを計算する
// そしてその各パターンから、各QuizLineItemの取りうる最大の範囲を計算する
func CalculateItemRangeListPatterns(quizItemAllocationPatterns quizpattern.QuizItemAllocationPatterns, quizLineLength int, splittedAnswerLine splitanswerline.SplittedAnswerLine) ItemRangeList {
	// 各QuizLineItemの割り振りパターンでのQuizLineItemの各入りうる範囲
	itemRangeListPatterns := make(ItemRangeListPatterns, len(quizItemAllocationPatterns))

	// QuizLineItemの割り振りパターンでforを回す
	for h, itemAllocationPattern := range quizItemAllocationPatterns {
		itemRangeList := make(ItemRangeList, 0, quizLineLength)
		// その中で、splitされた一箇所の、QuizLineItemとAnswerの値を取り出す
		for i := range splittedAnswerLine {
			// QuizLineの一部分。構造体なので、全体のQuizLineの何番目かのindexも取れる
			// [3,1,2]
			partQuizLineItem := itemAllocationPattern[i]
			// AnswerLineの一部分
			// [_,_,_,◼,_,_,_,_,_]
			answer := splittedAnswerLine[i]

			fmt.Printf("=====\n　確認中の問題の一部: %+v\n　その数値が入る解答欄の一部: %+v\n", partQuizLineItem, answer)

			filledOwnerIndexSlice := make([]int, len(answer.Cells))
			for index := range filledOwnerIndexSlice {
				filledOwnerIndexSlice[index] = -1
			}

			// QuizLineの数値ひとつに注目
			for j, quizItem := range partQuizLineItem {
				// answerから見た、そのquizItemが入りうる範囲
				start, end := partQuizLineItem.SidesLength(j)
				end = answer.Length() - 1 - end

				fmt.Println(start, end)

				for k := start; k <= end; k++ {
					// 範囲の中のFilledがそのquizItemだけのものか記録していく
					if answer.Cells[k] == schemas.Filled {
						// 誰もそのFilledの所有権を主張していなければ(-1)、quizItemのindexを入れる
						if filledOwnerIndexSlice[k] == -1 {
							filledOwnerIndexSlice[k] = quizItem.Index
						} else {
							// 既に誰か所有していれば(-1以外)、-2としてダメFilledとする
							filledOwnerIndexSlice[k] = -2
						}
					}
				}

				// start,endはグローバルな一次元座標なので、オフセットとなるanswer.startを加算している
				ir := ItemRange{
					start: answer.Start + start,
					end:   answer.Start + end,
					item:  quizItem,
				}
				// 一つ前のitemRangeのfilledEndよりもstartが左にある場合、手放す必要があるのでずらす
				if j != 0 && itemRangeList[len(itemRangeList)-1].filledEnd != nil {
					ir.start = max(ir.start, *itemRangeList[len(itemRangeList)-1].filledEnd+1)
				}

				fmt.Printf("ItemRange: %+v\n", ir)

				itemRangeList = append(itemRangeList, ir)
			}

			fmt.Printf("各数値が入る範囲: %+v\n", itemRangeList)

			fmt.Println(filledOwnerIndexSlice)
			for answerIndex, quizIndex := range filledOwnerIndexSlice {
				if quizIndex > -1 {
					fmt.Printf("調整前のitemRange: %+v\n", itemRangeList[quizIndex])
					itemRange := itemRangeList[quizIndex]
					length := itemRange.item.Value
					itemRangeList[quizIndex].start = max(itemRange.start, answer.Start+answerIndex-(length-1))
					itemRangeList[quizIndex].end = min(itemRange.end, answer.Start+answerIndex+(length-1))
					fmt.Printf("調整後のitemRange: %+v\n", itemRangeList[quizIndex])
				}
			}
		}
		// 1パターン分できたのでpush
		itemRangeListPatterns[h] = itemRangeList
	}
	fmt.Printf("各パターンごとの数値の取りうる範囲\n　　%+v\n", itemRangeListPatterns)

	// 各パターンを回し、各QuizLineItemの入りうる範囲の一番広い範囲を計算する
	maxItemRangeList := itemRangeListPatterns[0]
	for _, itemRangeList := range itemRangeListPatterns {
		for i, itemRange := range itemRangeList {
			maxItemRangeList[i].start = min(itemRange.start, maxItemRangeList[i].start)
			maxItemRangeList[i].end = max(itemRange.end, maxItemRangeList[i].end)
		}
	}
	return maxItemRangeList
}
