package handlers

import (
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
	expected := [][]schemas.CellType{
		{
			schemas.Filled,
			schemas.Unsettled,
			schemas.Unsettled,
			schemas.Unsettled,
		},
		{
			schemas.Unsettled,
			schemas.Filled,
			schemas.Unsettled,
			schemas.Unsettled,
		},
	}

	splittedAnswerLine := answerLine.splitAnswerLine()

	if !reflect.DeepEqual(splittedAnswerLine, expected) {
		t.Error()
	}
}
