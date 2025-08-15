package schemas

type Quiz struct {
	Title      string  `json:"title"`
	Answer     string  `json:"answer"`
	HandleName string  `json:"handlename"`
	Date       string  `json:"date"`
	Horizontal [][]int `json:"horz"`
	Vertical   [][]int `json:"vert"`
}
