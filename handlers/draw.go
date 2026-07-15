package handlers

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"picross/logger"
	"picross/schemas"
)

const (
	marginX         = 40
	marginTop       = 120
	marginBottom    = 30
	cellWidth       = 30
	cellBorderWidth = 2
)

var (
	backgroundColor = color.White
	borderColor     = color.Black
	filledColor     = color.RGBA{32, 32, 32, 255}
	blankColor      = color.RGBA{200, 200, 200, 255}
	unsettledColor  = color.White
)

// 回答用の配列を見て、指定されたパスに回答のイラストを出力する
func DrawAnswerImage(answer schemas.Answer, outputPath string) {
	horizontalCellNumber := answer.GetLength(schemas.Horizontal)
	verticalCellNumber := answer.GetLength(schemas.Vertical)

	// 大元の画像を作成
	img := image.NewRGBA(
		image.Rect(
			0,
			0,
			marginX*2+(cellBorderWidth+cellWidth)*horizontalCellNumber+cellBorderWidth,
			marginTop+marginBottom+(cellBorderWidth+cellWidth)*verticalCellNumber+cellBorderWidth,
		),
	)
	// 白で塗りつぶす
	draw.Draw(img, img.Bounds(), &image.Uniform{backgroundColor}, image.Point{0, 0}, draw.Src)
	// 回答を描画する矩形を黒で塗りつぶす
	draw.Draw(
		img,
		image.Rect(
			marginX,
			marginTop,
			marginX+(cellBorderWidth+cellWidth)*horizontalCellNumber+cellBorderWidth,
			marginTop+(cellBorderWidth+cellWidth)*verticalCellNumber+cellBorderWidth,
		),
		&image.Uniform{borderColor},
		image.Point{0, 0}, draw.Src,
	)

	answer.Map(func(x, y int, data schemas.CellType) {
		var cellColor color.Color
		switch data {
		case schemas.Filled:
			cellColor = filledColor
		case schemas.Unfilled:
			cellColor = blankColor
		case schemas.Unsettled:
			cellColor = unsettledColor
		}
		// マスを描画
		draw.Draw(
			img,
			image.Rect(
				marginX+cellBorderWidth+(cellWidth+cellBorderWidth)*x,
				marginTop+cellBorderWidth+(cellWidth+cellBorderWidth)*y,
				marginX+(cellBorderWidth+cellWidth)*(x+1),
				marginTop+(cellBorderWidth+cellWidth)*(y+1),
			),
			&image.Uniform{cellColor},
			image.Point{0, 0},
			draw.Src,
		)
	})

	// 既にファイルが存在している場合は削除する
	if _, err := os.Stat(outputPath); err == nil {
		os.Remove(outputPath)
	} else if !os.IsNotExist(err) {
		logger.ErrorLog(fmt.Sprintf("ファイルの存在確認中にエラーが発生しました: %+v", err))
		return
	}
	// ファイルを保存
	saveFile, _ := os.Create(outputPath)
	defer saveFile.Close()
	png.Encode(saveFile, img)
}
