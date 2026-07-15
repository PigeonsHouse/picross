package splitanswerline

import (
	"picross/schemas"
	"reflect"
	"testing"
)

func TestSplitAnswerLine(t *testing.T) {
	var answerLine = []schemas.CellType{
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
			Start: 0,
			End:   3,
			Cells: []schemas.CellType{
				schemas.Filled,
				schemas.Unsettled,
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
		SplittedAnswerLinePart{
			Start: 5,
			End:   8,
			Cells: []schemas.CellType{
				schemas.Unsettled,
				schemas.Filled,
				schemas.Unsettled,
				schemas.Unsettled,
			},
		},
	}

	got := SplitAnswerLine(answerLine)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Diff failed. \nExpected:\n%+v\n\nGot:\n%+v", expected, got)
	}
}
