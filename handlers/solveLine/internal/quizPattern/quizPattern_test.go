package quizpattern

import (
	splitanswerline "picross/handlers/solveLine/internal/splitAnswerLine"
	"picross/schemas"
	"reflect"
	"testing"
)

func TestGenerateQuizPatterns1(t *testing.T) {
	quizLine := []int{3, 1, 2}
	sal := splitanswerline.SplittedAnswerLine{
		splitanswerline.SplittedAnswerLinePart{
			Start: 0,
			End:   7,
			Cells: []schemas.CellType{
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
				QuizLineItem{Index: 0, Value: 3},
				QuizLineItem{Index: 1, Value: 1},
				QuizLineItem{Index: 2, Value: 2},
			},
		},
	}

	got := GenerateQuizPatterns(sal, quizLine)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Diff failed. \nExpected:\n%+v\n\nGot:\n%+v", expected, got)
	}
}

func TestGenerateQuizPatterns2(t *testing.T) {
	quizLine := []int{3, 1, 2}
	sal := splitanswerline.SplittedAnswerLine{
		splitanswerline.SplittedAnswerLinePart{
			Start: 0,
			End:   4,
			Cells: []schemas.CellType{
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
		splitanswerline.SplittedAnswerLinePart{
			Start: 5,
			End:   6,
			Cells: []schemas.CellType{
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
		splitanswerline.SplittedAnswerLinePart{
			Start: 7,
			End:   8,
			Cells: []schemas.CellType{
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
	}
	expected := QuizItemAllocationPatterns{
		QuizItemAllocationPattern{
			QuizItemAllocationInPart{
				QuizLineItem{Index: 0, Value: 3},
			},
			QuizItemAllocationInPart{
				QuizLineItem{Index: 1, Value: 1},
			},
			QuizItemAllocationInPart{
				QuizLineItem{Index: 2, Value: 2},
			},
		},
		QuizItemAllocationPattern{
			QuizItemAllocationInPart{
				QuizLineItem{Index: 0, Value: 3},
				QuizLineItem{Index: 1, Value: 1},
			},
			QuizItemAllocationInPart{},
			QuizItemAllocationInPart{
				QuizLineItem{Index: 2, Value: 2},
			},
		},
		QuizItemAllocationPattern{
			QuizItemAllocationInPart{
				QuizLineItem{Index: 0, Value: 3},
				QuizLineItem{Index: 1, Value: 1},
			},
			QuizItemAllocationInPart{
				QuizLineItem{Index: 2, Value: 2},
			},
			QuizItemAllocationInPart{},
		},
	}

	got := GenerateQuizPatterns(sal, quizLine)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Diff failed. \nExpected:\n%+v\n\nGot:\n%+v", expected, got)
	}
}
