package schemas

type Orientation int

const (
	// 水平方向
	Horizontal Orientation = iota
	// 垂直方向
	Vertical
)

func (o Orientation) String() string {
	switch o {
	case Horizontal:
		return "水平方向"
	case Vertical:
		return "垂直方向"
	default:
		panic("存在しないCellTypeです")
	}
}
