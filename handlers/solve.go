package handlers

import (
	"piclos/schemas"
)

func firstSolveLine(maxLength int, quizLine []int, answerLine []int) (isChanged bool) {
	/// quizLineの最小の長さを計算
	// 黒の数を数える
	sum := 0
	for _, v := range quizLine {
		sum += v
	}
	// 黒の間の最小スペースを数える
	sum += len(quizLine) - 1

	// 最小の長さとラインの長さの差
	lengthDiff := maxLength - sum
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

// ピクロスの問題の構造体を受け取り、回答の配列に値を入れていく
// 最後まで埋めることが出来たら、解けたかどうか示す返り値のフラグを立てる
func SolveQuiz(quiz schemas.Quiz, answer *schemas.Answer) bool {
	// 現在見ているラインが横向き(行)かどうか
	isHorizontal := true
	// 横向きのラインが全て終わったかどうか
	isHorizontalFinish := false
	// 縦向きのラインが全て終わったかどうか
	isVerticalFinish := false

	for {
		isChanged := false
		// 向きを切り替える
		isHorizontal = !isHorizontal

		for lineNum := 0; lineNum < answer.GetLength(isHorizontal); lineNum++ {
			var (
				// 見ているラインの問題データ (ex. [3, 2, 1, 1, 4])
				quizLine []int
				// 見ているラインの回答データ (ex. [1, 0, 1, -1, 0, ..., 1])
				answerLine []int
			)

			// 問題のラインと回答用のラインを取得
			if isHorizontal {
				quizLine = quiz.Horizontal[lineNum]
			} else {
				quizLine = quiz.Vertical[lineNum]
			}
			answer.CopyLine(isHorizontal, lineNum, &answerLine)

			// answerLineを埋める関数を列挙していく
			isChanged = isChanged || firstSolveLine(answer.GetLength(isHorizontal), quizLine, answerLine)
			// TODO: 他の解法関数を追加する
			// 例) isChanged = isChanged || someSolveFunction()

			// 変化があった場合は、answerに保存する
			if isChanged {
				answer.SaveLine(isHorizontal, lineNum, answerLine)
			}
		}

		// 変化がなければ、見てる向きのラインが全て終わったとする
		if !isChanged {
			if isHorizontal {
				isHorizontalFinish = true
			} else {
				isVerticalFinish = true
			}
		} else {
			// 変化があったので、両方終わってない判定に戻す
			isHorizontalFinish = false
			isVerticalFinish = false
		}
		// 両方の向きで全ラインに変化がなければ終了
		if isHorizontalFinish && isVerticalFinish {
			break
		}
	}

	return answer.IsSolved()
}
