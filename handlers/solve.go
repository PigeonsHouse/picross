package handlers

import (
	"piclos/schemas"
)

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

			isChanged = isChanged || SolveLine(quizLine, answerLine)

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
