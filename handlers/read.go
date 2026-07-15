package handlers

import (
	"encoding/json"
	"os"
	"picross/logger"
	"picross/schemas"
)

// ピクロスの問題ファイルを受け取り、データをパースして、問題のサイズにあった回答用の配列を作成する
func ReadQuiz(fileName string) (quiz schemas.Quiz, answer schemas.Answer) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		logger.ErrorLog(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, &quiz)

	answer.Initialize(quiz.GetSize())

	return
}
