package path

import vec "github.com/kaneta1992/vector/vector2"

type Glyph struct {
	Contours []*Contour
}

func (g *Glyph) CreatePointsAndInnerSegments() ([]*vec.Vector2, [][2]int32, []*vec.Vector2) {
	retPts := []*vec.Vector2{}
	retSegs := [][2]int32{}
	retHoles := []*vec.Vector2{}
	for _, c := range g.Contours {
		vertexNum := len(retPts)
		points := c.getPoints()
		retPts = append(retPts, points...)

		segments := c.CreateInnerSegments(vertexNum)
		retSegs = append(retSegs, segments...)

		hole := c.getHolePoint()
		if hole != nil {
			retHoles = append(retHoles, hole)
		}
	}
	return retPts, retSegs, retHoles
}

func (g *Glyph) AddContour(c *Contour) {
	g.Contours = append(g.Contours, c)
}
