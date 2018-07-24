package path

import (
	"github.com/kaneta1992/vector/vector2"
)

type ISegment interface {
	getPoints() []*vec.Vector2
}

type ICurve interface {
	splitCurve() *Curves
}
