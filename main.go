package main

import (
	"flag"
	"os"
	"picross/handlers"
	"picross/logger"
)

func main() {
	logLevel := logger.Progress
	logger.Init(logLevel)

	outputPath := flag.String("o", "answer.png", "quiz file path")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		logger.ErrorLog("コマンド引数で、解くピクロスのjsonファイルのパスを指定してください")
		os.Exit(1)
	}
	inputPath := args[0]

	quiz, ans := handlers.ReadQuiz(inputPath)
	if isSolved := handlers.SolveQuiz(quiz, &ans); isSolved {
		handlers.DrawAnswerImage(ans, *outputPath)
	} else {
		logger.InfoLog("与えられたピクロスが解けませんでした")
		if logLevel == logger.Debug {
			handlers.DrawAnswerImage(ans, *outputPath)
		}
	}
}
