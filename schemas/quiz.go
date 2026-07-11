package schemas

type Quiz struct {
	Title      string  `json:"title"`
	Answer     string  `json:"answer"`
	HandleName string  `json:"handlename"`
	Date       string  `json:"date"`
	Horizontal [][]int `json:"horz"`
	Vertical   [][]int `json:"vert"`
}

func (q Quiz) GetSize() (int, int) {
	return len(q.Horizontal), len(q.Vertical)
}

func (q Quiz) ReadLine(orientation Orientation, index int) []int {
	if orientation == Horizontal {
		return q.Horizontal[index]
	} else {
		return q.Vertical[index]
	}
}
