package path

import vec "github.com/kaneta1992/vector/vector2"

type Glyph struct {
	Contours []*Contour
}

func (g *Glyph) CreatePointsAndInnerSegments() ([]*vec.Vector2, [][2]int32) {
	retPts := []*vec.Vector2{}
	retSegs := [][2]int32{}
	for _, c := range g.Contours {
		vertexNum := len(retPts)
		points := c.getPoints()
		segments := c.CreateInnerSegments(vertexNum)
		retPts = append(retPts, points...)
		retSegs = append(retSegs, segments...)
	}
	return retPts, retSegs
}

func (g *Glyph) AddContour(c *Contour) {
	g.Contours = append(g.Contours, c)
}
