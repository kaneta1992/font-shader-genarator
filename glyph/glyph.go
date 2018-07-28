package path

import vec "github.com/kaneta1992/vector/vector2"

type Glyph struct {
	Contours []*Contour
}

func (g *Glyph) CreatePointsAndInnerSegments() ([]*vec.Vector2, [][2]int32, []*vec.Vector2, [][3]int32) {
	retPts := []*vec.Vector2{}
	retSegs := [][2]int32{}
	retHoles := []*vec.Vector2{{999999.9, 999999.9}}
	retBeziers := [][3]int32{}
	for _, c := range g.Contours {
		vertexNum := len(retPts)
		points := c.getPoints()
		retPts = append(retPts, points...)

		segments, beziers := c.CreateInnerSegments(vertexNum)
		retSegs = append(retSegs, segments...)
		retBeziers = append(retBeziers, beziers...)

		hole := c.getHolePoint()
		if hole != nil {
			retHoles = append(retHoles, hole)
		}
	}
	return retPts, retSegs, retHoles, retBeziers
}

func (g *Glyph) AddContour(c *Contour) {
	g.Contours = append(g.Contours, c)
}
