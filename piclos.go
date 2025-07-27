package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
)

type Piclos struct {
	Title      string  `json:"title"`
	Answer     string  `json:"answer"`
	HandleName string  `json:"handlename"`
	Date       string  `json:"date"`
	Horizontal [][]int `json:"horz"`
	Vertical   [][]int `json:"vert"`
}

func ReadQuiz(fileName string, quiz *Piclos) (answer [][]int) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, quiz)

	hor := len(quiz.Horizontal)
	ver := len(quiz.Vertical)

	answer = make([][]int, ver)
	for i := range answer {
		answer[i] = make([]int, hor)
	}
	return
}

func getLine(horizontal bool, num int, ans [][]int) (line []int) {
	if horizontal {
		line = make([]int, len(ans))
		for i, v := range ans[num] {
			line[i] = v
		}
	} else {
		line = make([]int, len(ans))
		for i, l := range ans {
			line[i] = l[num]
		}
	}
	return
}

func saveLine(horizontal bool, num int, ans [][]int, line []int) {
	if horizontal {
		for i, v := range line {
			ans[num][i] = v
		}
	} else {
		for i, v := range line {
			ans[i][num] = v
		}
	}
	return
}

func SolvingQuiz(quiz Piclos, answer [][]int) (isSolved bool) {
	var (
		isHorizontal bool
		lineNum      int
		line         []int
		quizLine     []int
	)

	for i := 0; i < len(answer); i++ {
		isHorizontal = true
		lineNum = i
		line = getLine(isHorizontal, lineNum, answer)

		quizLine = quiz.Horizontal[i]

		sum := 0
		for _, v := range quizLine {
			sum += v
		}
		sum += len(quizLine) - 1

		if len(answer[0])-sum < 0 {
			continue
		}

		line[2] = 1

		saveLine(isHorizontal, lineNum, answer, line)
	}

	for {
		// ここでどのラインを考慮するか考える
		isHorizontal = true
		lineNum = 1

		line = getLine(isHorizontal, lineNum, answer)

		// ここでラインの中で黒(1)やバツ(-1)を決める
		line[2] = 1

		// 終わったら保存
		saveLine(isHorizontal, lineNum, answer, line)

		break
	}
	return
}

func GenerateAnswerImage(quiz Piclos, answer [][]int, isSolved bool, outputPath string) {
	my := 40
	mt := 120
	mb := 30
	bw := 30
	bf := 2
	if _, err := os.Stat(outputPath); err == nil {
		os.Remove(outputPath)
	}

	img := image.NewRGBA(image.Rect(0, 0, 2*my+bf+bw*len(answer[0]), mt+mb+bf+bw*len(answer)))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{0, 0}, draw.Src)
	draw.Draw(img, image.Rect(my, mt, my+bf+bw*len(answer[0]), mt+bf+bw*len(answer)), &image.Uniform{color.Black}, image.Point{0, 0}, draw.Src)

	for y, inner := range answer {
		for x, data := range inner {
			var fillColor color.Color
			if data == 1 {
				fillColor = color.RGBA{32, 32, 32, 255}
			} else if data == -1 {
				fillColor = color.RGBA{200, 200, 200, 255}
			} else {
				fillColor = color.White
			}
			fmt.Println(data)
			draw.Draw(img, image.Rect(my+bw*x+bf, mt+bw*y+bf, my+bw*(x+1), mt+bw*(y+1)), &image.Uniform{fillColor}, image.Point{0, 0}, draw.Src)
		}
	}

	savefile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer savefile.Close()
	png.Encode(savefile, img)
	return
}

func main() {
	outPath := flag.String("o", "answer.png", "quiz file path")
	flag.Parse()
	args := flag.Args()

	quiz := Piclos{}
	ans := ReadQuiz(args[0], &quiz)

	isSolved := SolvingQuiz(quiz, ans)
	GenerateAnswerImage(quiz, ans, isSolved, *outPath)
}
