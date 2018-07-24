package path

import vec "github.com/kaneta1992/vector/vector2"

type Curves struct {
	Curves []*Curve
}

// ベジエ曲線を分割する
func (c *Curves) splitCurve() *Curves {
	return nil
}

func (c *Curves) getPoints() []*vec.Vector2 {
	return []*vec.Vector2{}
}
