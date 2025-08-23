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
	marginX := 40
	marginTop := 120
	marginBottom := 30
	boxWidth := 30
	boxBorderWidth := 2
	horizontalCellNumber := answer.GetLength(true)
	verticalCellNumber := answer.GetLength(false)

	// 大元の画像を作成
	img := image.NewRGBA(
		image.Rect(
			0,
			0,
			marginX*2+(boxBorderWidth+boxWidth)*horizontalCellNumber+boxBorderWidth,
			marginTop+marginBottom+(boxBorderWidth+boxWidth)*verticalCellNumber+boxBorderWidth,
		),
	)
	// 白で塗りつぶす
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{0, 0}, draw.Src)
	// 回答を描画する矩形を黒で塗りつぶす
	draw.Draw(
		img,
		image.Rect(
			marginX,
			marginTop,
			marginX+(boxBorderWidth+boxWidth)*horizontalCellNumber+boxBorderWidth,
			marginTop+(boxBorderWidth+boxWidth)*verticalCellNumber+boxBorderWidth,
		),
		&image.Uniform{color.Black},
		image.Point{0, 0}, draw.Src,
	)

	for y, inner := range answer.GetData() {
		for x, data := range inner {
			// マスの中身に応じて色を変える
			// 1: 黒、-1: 灰色、0: 白
			var fillColor color.Color
			if data == 1 {
				fillColor = color.RGBA{32, 32, 32, 255}
			} else if data == -1 {
				fillColor = color.RGBA{200, 200, 200, 255}
			} else {
				fillColor = color.White
			}
			// マスを描画
			draw.Draw(
				img,
				image.Rect(
					marginX+boxBorderWidth+(boxWidth+boxBorderWidth)*x,
					marginTop+boxBorderWidth+(boxWidth+boxBorderWidth)*y,
					marginX+(boxBorderWidth+boxWidth)*(x+1),
					marginTop+(boxBorderWidth+boxWidth)*(y+1),
				),
				&image.Uniform{fillColor},
				image.Point{0, 0},
				draw.Src,
			)
		}
	}

	// 既にファイルが存在している場合は削除する
	if _, err := os.Stat(outputPath); err == nil {
		os.Remove(outputPath)
	} else if !os.IsNotExist(err) {
		fmt.Println("ファイルの存在確認中にエラーが発生しました:", err)
		return
	}
	// ファイルを保存
	saveFile, _ := os.Create(outputPath)
	defer saveFile.Close()
	png.Encode(saveFile, img)
}
