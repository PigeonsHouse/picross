package handlers

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"piclos/schemas"
)

// 回答用の配列を見て、指定されたパスに回答のイラストを出力する
func DrawAnswerImage(answer schemas.Answer, outputPath string) {
	my := 40
	mt := 120
	mb := 30
	bw := 30
	bf := 2
	horizontalLength := answer.GetLength(true)
	verticalLength := answer.GetLength(false)

	// 大元の画像を作成
	img := image.NewRGBA(
		image.Rect(
			0,
			0,
			2*my+bf+bw*horizontalLength,
			mt+mb+bf+bw*verticalLength,
		),
	)
	// 白で塗りつぶす
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{0, 0}, draw.Src)
	// 回答を描画する矩形を黒で塗りつぶす
	draw.Draw(
		img,
		image.Rect(
			my,
			mt,
			my+bf+bw*horizontalLength,
			mt+bf+bw*verticalLength,
		),
		&image.Uniform{color.Black},
		image.Point{0, 0}, draw.Src,
	)

	for y, inner := range answer.GetData() {
		for x, data := range inner {
			var fillColor color.Color
			if data == 1 {
				fillColor = color.RGBA{32, 32, 32, 255}
			} else if data == -1 {
				fillColor = color.RGBA{200, 200, 200, 255}
			} else {
				fillColor = color.White
			}
			draw.Draw(
				img,
				image.Rect(
					my+bw*x+bf,
					mt+bw*y+bf,
					my+bw*(x+1),
					mt+bw*(y+1),
				),
				&image.Uniform{fillColor},
				image.Point{0, 0},
				draw.Src,
			)
		}
	}

	if _, err := os.Stat(outputPath); err == nil {
		os.Remove(outputPath)
	} else if !os.IsNotExist(err) {
		fmt.Println("ファイルの存在確認中にエラーが発生しました:", err)
		return
	}
	saveFile, _ := os.Create(outputPath)
	defer saveFile.Close()
	png.Encode(saveFile, img)
}
