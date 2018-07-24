package path

import vec "github.com/kaneta1992/vector/vector2"

type Line struct {
	Start *vec.Vector2
	End   *vec.Vector2
}

func (l *Line) getPoints() []*vec.Vector2 {
	return []*vec.Vector2{l.End}
}
