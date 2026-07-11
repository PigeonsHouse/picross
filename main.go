package main

import (
	"flag"
	"fmt"
	"os"
	"picross/handlers"
)

func main() {
	outputPath := flag.String("o", "answer.png", "quiz file path")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("[ERR] コマンド引数で、解くピクロスのjsonファイルのパスを指定してください")
		os.Exit(1)
	}
	inputPath := args[0]

	quiz, ans := handlers.ReadQuiz(inputPath)
	if isSolved := handlers.SolveQuiz(quiz, &ans); isSolved {
		handlers.DrawAnswerImage(ans, *outputPath)
	} else {
		fmt.Println("与えられたピクロスが解けませんでした")
		// デバッグのために、解けなかった場合も画像を出力する
		handlers.DrawAnswerImage(ans, *outputPath)
	}
}
