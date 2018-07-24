package path

import vec "github.com/kaneta1992/vector/vector2"

type Point struct {
	Point *vec.Vector2
}

func (m *Point) getPoints() []*vec.Vector2 {
	return []*vec.Vector2{m.Point}
}
