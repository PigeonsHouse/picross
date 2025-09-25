package handlers

import (
	"reflect"
	"testing"
)

func TestSplitAnswerLine(t *testing.T) {
	answerLine := []int{1, 0, 0, 0, -1, 0, 1, 0, 0}
	expected := [][]int{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
	}

	splittedAnswerLine := SplitAnswerLine(answerLine)

	if !reflect.DeepEqual(splittedAnswerLine, expected) {
		t.Error()
	}
}
