package handlers

import (
	"encoding/json"
	"picross/schemas"
	"reflect"
	"testing"
)

func TestSplitAnswerLine(t *testing.T) {
	var answerLine AnswerLine = []schemas.CellType{
		schemas.Filled,
		schemas.Unsettled,
		schemas.Unsettled,
		schemas.Unsettled,
		schemas.Unfilled,
		schemas.Unsettled,
		schemas.Filled,
		schemas.Unsettled,
		schemas.Unsettled,
	}
	expected := SplittedAnswerLine{
		SplittedAnswerLinePart{
			start: 0,
			end:   3,
			cells: []schemas.CellType{
				schemas.Filled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
		SplittedAnswerLinePart{
			start: 4,
			end:   7,
			cells: []schemas.CellType{
				schemas.Unsettled,
				schemas.Filled,
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
	}

	splittedAnswerLine := answerLine.splitAnswerLine()

	if !reflect.DeepEqual(splittedAnswerLine, expected) {
		expBytes, _ := json.MarshalIndent(expected, "", "  ")
		gotBytes, _ := json.MarshalIndent(splittedAnswerLine, "", "  ")

		t.Errorf("Diff failed. \nExpected:\n%s\n\nGot:\n%s", string(expBytes), string(gotBytes))
	}
}

func TestGenerateQuizPatterns1(t *testing.T) {
	quizLine := []int{3, 1, 2}
	sal := SplittedAnswerLine{
		SplittedAnswerLinePart{
			start: 0,
			end:   7,
			cells: []schemas.CellType{
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
	}
	expected := QuizItemAllocationPatterns{
		QuizItemAllocationPattern{
			QuizItemAllocationInPart{
				QuizLineItem{index: 0, value: 3},
				QuizLineItem{index: 1, value: 1},
				QuizLineItem{index: 2, value: 2},
			},
		},
	}

	quizPatterns := sal.generateQuizPatterns(quizLine)

	if !reflect.DeepEqual(quizPatterns, expected) {
		expBytes, _ := json.MarshalIndent(expected, "", "  ")
		gotBytes, _ := json.MarshalIndent(quizPatterns, "", "  ")

		t.Errorf("Diff failed. \nExpected:\n%s\n\nGot:\n%s", string(expBytes), string(gotBytes))
	}
}

func TestGenerateQuizPatterns2(t *testing.T) {
	quizLine := []int{3, 1, 2}
	sal := SplittedAnswerLine{
		SplittedAnswerLinePart{
			start: 0,
			end:   4,
			cells: []schemas.CellType{
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
		SplittedAnswerLinePart{
			start: 4,
			end:   6,
			cells: []schemas.CellType{
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
		SplittedAnswerLinePart{
			start: 6,
			end:   8,
			cells: []schemas.CellType{
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
	}
	expected := QuizItemAllocationPatterns{
		QuizItemAllocationPattern{
			QuizItemAllocationInPart{
				QuizLineItem{index: 0, value: 3},
			},
			QuizItemAllocationInPart{
				QuizLineItem{index: 1, value: 1},
			},
			QuizItemAllocationInPart{
				QuizLineItem{index: 2, value: 2},
			},
		},
		QuizItemAllocationPattern{
			QuizItemAllocationInPart{
				QuizLineItem{index: 0, value: 3},
				QuizLineItem{index: 1, value: 1},
			},
			QuizItemAllocationInPart{
				QuizLineItem{index: 2, value: 2},
			},
			QuizItemAllocationInPart{},
		},
		QuizItemAllocationPattern{
			QuizItemAllocationInPart{
				QuizLineItem{index: 0, value: 3},
				QuizLineItem{index: 1, value: 1},
			},
			QuizItemAllocationInPart{},
			QuizItemAllocationInPart{
				QuizLineItem{index: 2, value: 2},
			},
		},
	}

	quizPatterns := sal.generateQuizPatterns(quizLine)

	if !reflect.DeepEqual(quizPatterns, expected) {
		expBytes, _ := json.MarshalIndent(expected, "", "  ")
		gotBytes, _ := json.MarshalIndent(quizPatterns, "", "  ")

		t.Errorf("Diff failed. \nExpected:\n%s\n\nGot:\n%s", string(expBytes), string(gotBytes))
	}
}
