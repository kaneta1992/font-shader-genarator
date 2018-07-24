package path

import (
	"github.com/kaneta1992/vector/vector2"
)

type Contour struct {
	Segments []ISegment
	nowPoint *vec.Vector2
}

func (c *Contour) CreateInnerLines() ([]*vec.Vector2, []int32) {
	return nil, nil
}

func (c *Contour) GenerateShaderCode() string {
	return ""
}

func (c *Contour) addSegment(seg ISegment) {
	c.Segments = append(c.Segments, seg)
}

func (c *Contour) ToMove(start *vec.Vector2) {
	c.addSegment(&Point{start})
	c.nowPoint = start
}

func (c *Contour) ToLine(end *vec.Vector2) {
	c.addSegment(&Line{c.nowPoint, end})
	c.nowPoint = end
}

func (c *Contour) ToCurve(control, end *vec.Vector2) {
	c.addSegment(&Curve{c.nowPoint, control, end})
	c.nowPoint = end
}

func (c *Contour) NumPoints() int {
	return len(c.getPoints())
}

func wrap(x, y int) int {
	return ((x % y) + y) % y
}

func (c *Contour) getPoints() []*vec.Vector2 {
	ret := []*vec.Vector2{}
	for _, s := range c.Segments {
		ret = append(ret, s.getPoints()...)
	}
	return ret
}

func (c *Contour) CreateInnerSegments(adjust int) [][2]int32 {
	ret := [][2]int32{}
	offset := 0
	max := c.NumPoints()
	for _, s := range c.Segments {
		switch s.(type) {
		case *Line:
			ret = append(ret, [][2]int32{{int32(wrap(offset-1, max) + adjust), int32(offset + adjust)}}...)
			offset++
		case ICurve:
			ret = append(ret, [][2]int32{{int32(wrap(offset-1, max) + adjust), int32(offset + adjust)}, {int32(offset + adjust), int32(offset + 1 + adjust)}}...)
			offset += 2
		}
	}
	ret = append(ret, [][2]int32{{int32(offset + adjust - 1), int32(offset + adjust)}}...)
	return ret
}
