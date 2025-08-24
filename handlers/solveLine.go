package handlers

import "reflect"

// answerLineを埋める関数をまとめる
func SolveLine(quizLine []int, answerLine []int) (isChanged bool) {
	if len(quizLine) == 0 {
		panic("quizLineの長さが0")
	}
	// answerLineに0が含まれていない場合は、解き終わっているので何もしない
	if !hasZero(answerLine) {
		return
	}
	isChanged = isChanged || withoutEdgeCross(solveLineFirst, quizLine, answerLine)
	isChanged = isChanged || withoutEdgeCross(solveLineEdge, quizLine, answerLine)
	// TODO: 他の解法関数を追加する
	// 例) isChanged = isChanged || someSolveFunction()

	isChanged = isChanged || fillBrankInComplete(quizLine, answerLine)

	return
}

func hasZero(answerLine []int) bool {
	for _, v := range answerLine {
		if v == 0 {
			return true
		}
	}
	return false
}

func withoutEdgeCross(fn func(quizLine []int, answerLine []int) bool, quizLine []int, answerLine []int) (isChanged bool) {
	// 左端、右端に✕が確定している場合、その分を除外して関数を実行する
	// hasZeroで0が含まれていることは確認している
	leftOffset := -1
	rightOffset := -1
	for i, v := range answerLine {
		if v != -1 {
			if leftOffset == -1 {
				leftOffset = i
			}
			rightOffset = i + 1
		}
	}

	isChanged = isChanged || fn(quizLine, answerLine[leftOffset:rightOffset])

	return
}

func solveLineFirst(quizLine []int, answerLine []int) (isChanged bool) {
	maxLineLength := len(answerLine)

	/// quizLineの最小の長さを計算
	// 黒の数を数える
	sum := 0
	for _, v := range quizLine {
		sum += v
	}
	// 黒の間の最小スペースを数える
	sum += len(quizLine) - 1

	// 最小の長さとラインの長さの差
	lengthDiff := maxLineLength - sum
	// ラインの塗りが確定してるか
	isConfirmed := lengthDiff == 0
	// 最小の塗りの長さはラインの長さを超えるはずがない
	if lengthDiff < 0 {
		panic("quizの横ラインの合計の長さが答えの横ラインの長さより大きい")
	}

	avoidLength := lengthDiff
	blackLength := quizLine[0] - avoidLength
	// index=0でblackLengthを取得したので、次取るindexは1になる
	quizLineIndex := 1
	for i := range answerLine {
		// i番目の値
		value := 0
		if avoidLength > 0 {
			// avoidLengthが残っている場合は、0になる(0なのでvalueは変わらない)
			avoidLength--
		} else if blackLength > 0 {
			// blackLengthが残っている場合は、1になる
			value = 1
			blackLength--
		} else {
			// avoidLength、blackLengthの両方が0になったら、i番目は空白になる
			// 確定している場合は-1を入れる
			if isConfirmed {
				value = -1
			}
			// 空白を入れたので、次のavoidLengthとblackLengthを更新する
			if quizLineIndex < len(quizLine) {
				avoidLength = lengthDiff
				blackLength = quizLine[quizLineIndex] - avoidLength
				quizLineIndex++
			}
		}
		// 新しく変化があれば、更新してフラグを立てる
		if answerLine[i] == 0 && value != 0 {
			answerLine[i] = value
			isChanged = true
		}
	}
	return
}

func solveLineEdge(quizLine []int, answerLine []int) (isChanged bool) {
	// 左端が黒の場合、quizLineの最初の値分だけ1にする
	if answerLine[0] == 1 {
		offset := quizLine[0]
		for i := range answerLine {
			if i >= offset {
				if answerLine[i] == 0 {
					answerLine[i] = -1
					isChanged = true
				}
				break
			}
			if answerLine[i] == 0 {
				answerLine[i] = 1
				isChanged = true
			}
		}
	}
	// 右端が黒の場合、quizLineの最後の値分だけ1にする
	if answerLine[len(answerLine)-1] == 1 {
		offset := quizLine[len(quizLine)-1]
		for i := len(answerLine) - 1; i >= 0; i-- {
			if i <= len(answerLine)-1-offset {
				if answerLine[i] == 0 {
					answerLine[i] = -1
					isChanged = true
				}
				break
			}
			if answerLine[i] == 0 {
				answerLine[i] = 1
				isChanged = true
			}
		}
	}

	return
}

func fillBrankInComplete(quizLine []int, answerLine []int) (isChanged bool) {
	// answerLineが解き終わっているか確認する
	var currentLine []int
	valContinue := false
	for _, v := range answerLine {
		if v == 1 {
			if valContinue {
				currentLine[len(currentLine)-1] += 1
			} else {
				currentLine = append(currentLine, 1)
			}
		}
		valContinue = v == 1
	}
	// currentLineが空の場合、[0]とする
	if len(currentLine) == 0 {
		currentLine = []int{0}
	}
	// まだ解き終わっていない場合は何もしない
	if !reflect.DeepEqual(quizLine, currentLine) {
		return
	}
	// 0を-1に変える
	for i, v := range answerLine {
		if v == 0 {
			answerLine[i] = -1
			isChanged = true
		}
	}

	return
}
