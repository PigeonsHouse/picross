package quizpattern

import (
	"fmt"
	splitanswerline "picross/handlers/internal/solveLine/internal/splitAnswerLine"
	"picross/logger"
)

// QuizLineの数字一つ分のデータ
type QuizLineItem struct {
	// QuizLine内のIndex
	Index int
	// QuizLine内の数値
	Value int
}

func convertQuizLineItemList(quizLine []int) []QuizLineItem {
	quizLineItems := make([]QuizLineItem, len(quizLine))
	for i, num := range quizLine {
		quizLineItems[i] = QuizLineItem{
			Index: i,
			Value: num,
		}
	}
	return quizLineItems
}
func convertPrimitiveQuizLine(quizLineItems []QuizLineItem) []int {
	quizLine := make([]int, len(quizLineItems))
	for i, item := range quizLineItems {
		quizLine[i] = item.Value
	}
	return quizLine
}

// SplittedAnswerLineにQuizLineのItemをどう置くかというパターンの一部分
// Lengthは、最大でQuizLineの長さ、最小は0
type QuizItemAllocationInPart []QuizLineItem

func (qiaip QuizItemAllocationInPart) SidesLength(index int) (int, int) {
	start := 0
	for _, item := range qiaip[:index] {
		start += item.Value + 1
	}
	end := 0
	for _, item := range qiaip[index+1:] {
		end += item.Value + 1
	}
	return start, end
}

// SplittedAnswerLineにQuizLineのItemをどう置くかというパターン
// 常にSplittedAnswerLineの長さと一致する
type QuizItemAllocationPattern []QuizItemAllocationInPart

// 上のパターンのリスト
type QuizItemAllocationPatterns []QuizItemAllocationPattern

func GenerateQuizPatterns(sal splitanswerline.SplittedAnswerLine, quizLine []int) QuizItemAllocationPatterns {
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

	var recFunc func(innerQuizLine []QuizLineItem, innerSal splitanswerline.SplittedAnswerLine) QuizItemAllocationPatterns
	recFunc = func(innerQuizLine []QuizLineItem, innerSal splitanswerline.SplittedAnswerLine) QuizItemAllocationPatterns {
		logger.DebugLog(fmt.Sprintf("start recFunc\nquiz: %+v target answer: %+v", innerQuizLine, innerSal[0]))
		var wholePatterns QuizItemAllocationPatterns
		if len(innerQuizLine) == 0 {
			// innerSalの長さ分partを含める必要がある
			wholePatterns = append(wholePatterns,
				make(QuizItemAllocationPattern, len(innerSal)),
			)
		} else if len(innerSal) == 1 {
			if innerSal[0].IsContainable(convertPrimitiveQuizLine(innerQuizLine)) {
				wholePatterns = append(wholePatterns,
					QuizItemAllocationPattern{
						QuizItemAllocationInPart(innerQuizLine),
					},
				)
			} else {
				// 入り切らないケースはnilとしてしまう
				logger.DebugLog(fmt.Sprintf(
					"finish recFunc\nquiz: %+v target answer: %+v result: nil",
					innerQuizLine, innerSal[0],
				))
				return nil
			}
		} else {
			logger.DebugLog("current part: []")
			patterns := recFunc(innerQuizLine, innerSal[1:])
			if patterns != nil {
				for i := range patterns {
					patterns[i] = append(
						QuizItemAllocationPattern{QuizItemAllocationInPart{}},
						patterns[i]...,
					)
				}
				logger.DebugLog(fmt.Sprintf("append pattern: %+v", patterns))
				wholePatterns = append(wholePatterns, patterns...)
			}

			for i := range innerQuizLine {
				target := innerQuizLine[:i+1]
				if !innerSal[0].IsContainable(convertPrimitiveQuizLine(target)) {
					break
				}
				logger.DebugLog(fmt.Sprintf("current part: %+v", innerQuizLine[:i+1]))
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
					logger.DebugLog(fmt.Sprintf("append pattern: %+v", patterns))
					wholePatterns = append(wholePatterns, patterns...)
				}
			}
		}
		logger.DebugLog(fmt.Sprintf("finish recFunc\nquiz: %+v target answer: %+v result: %+v", innerQuizLine, innerSal[0], wholePatterns))
		return wholePatterns
	}
	return recFunc(workingQuizLineItems, sal)
}
