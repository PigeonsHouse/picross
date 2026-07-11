package handlers

import (
	"fmt"
	"picross/schemas"
	"time"
)

// ピクロスの問題の構造体を受け取り、回答の配列に値を入れていく
// 最後まで埋めることが出来たら、解けたかどうか示す返り値のフラグを立てる
func SolveQuiz(quiz schemas.Quiz, answer *schemas.Answer) bool {
	// 現在見ているラインが横向き(行)かどうか
	currentOrientation := schemas.Horizontal
	// その向きのラインが全て終わったかどうか
	isFinish := map[schemas.Orientation]bool{
		schemas.Horizontal: false,
		schemas.Vertical:   false,
	}

	for {
		isChanged := false
		for lineIndex := 0; lineIndex < answer.GetLength(currentOrientation); lineIndex++ {
			// 問題のラインと回答用のラインを取得
			quizLine := quiz.ReadLine(currentOrientation, lineIndex)
			var answerLine AnswerLine = answer.ReadLine(currentOrientation, lineIndex)
			fmt.Println("[start] orientation", currentOrientation, "index", lineIndex, "quizLine", quizLine, "answerLine", answerLine)

			// 変化があった場合は、answerに保存する
			if isChangedLine := answerLine.SolveLine(quizLine); isChangedLine {
				answer.WriteLine(currentOrientation, lineIndex, answerLine)
				isChanged = true
			}
			// debug
			fmt.Println("[end] orientation", currentOrientation, "index", lineIndex, "quizLine", quizLine, "answerLine", answerLine)
			fmt.Println()
			time.Sleep(100 * time.Millisecond)
		}

		// 変化がなければ、見てる向きのラインが全て終わったとする
		if !isChanged {
			isFinish[currentOrientation] = true
		} else {
			// 変化があったので、両方終わってない判定に戻す
			isFinish[schemas.Horizontal] = false
			isFinish[schemas.Vertical] = false
		}
		// 両方の向きで全ラインに変化がなければ終了
		if isFinish[schemas.Horizontal] && isFinish[schemas.Vertical] {
			break
		}

		currentOrientation = switchOrientation(currentOrientation)
	}

	return answer.IsSolved()
}

func switchOrientation(orientation schemas.Orientation) schemas.Orientation {
	if orientation == schemas.Horizontal {
		return schemas.Vertical
	} else {
		return schemas.Horizontal
	}
}
