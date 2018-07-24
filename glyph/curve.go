package path

import (
	"github.com/kaneta1992/vector/vector2"
)

type Curve struct {
	Start   *vec.Vector2
	Control *vec.Vector2
	End     *vec.Vector2
}

// ベジエ曲線を分割する
func (c *Curve) splitCurve() *Curves {
	return nil
}

func (c *Curve) getPoints() []*vec.Vector2 {
	return []*vec.Vector2{c.Control, c.End}
}
