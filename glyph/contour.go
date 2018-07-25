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

func signedArea(points []*vec.Vector2) float64 {
	lp := len(points)
	area := 0.0
	for i, _ := range points {
		v1 := points[i]
		v2 := points[wrap(i+1, lp)]
		area += v1.Cross(v2)
	}
	return area
}
func (c *Contour) CreateInnerSegments(offset int) [][2]int32 {
	ret := [][2]int32{}
	max := c.NumPoints()
	index := 0
	for _, s := range c.Segments {
		switch s.(type) {
		case *Point:
			start := wrap(index-1, max)
			end := index
			ret = append(ret, [][2]int32{{int32(start + offset), int32(end + offset)}}...)
			index = end
		case *Line:
			start := index
			end := index + 1
			ret = append(ret, [][2]int32{{int32(start + offset), int32(end + offset)}}...)
			index = end
		case *Curve:
			curve := s.(*Curve)
			start := index
			control := index + 1
			end := index + 2
			area := signedArea([]*vec.Vector2{curve.Start, curve.Control, curve.End})
			if area < 0.0 {
				// 左回りなら制御点が内部(右)にいるので制御点も含めた頂点を追加する
				ret = append(ret, [][2]int32{{int32(start + offset), int32(control + offset)}, {int32(control + offset), int32(end + offset)}}...)
			} else {
				// 右回りなら制御点が外部(左)にいるので制御点を無視する
				//points = append(points, v2.MulScalar(irate))
				ret = append(ret, [2]int32{int32(start + offset), int32(end + offset)})
			}
			index = end
		}
	}
	return ret
}
