package path

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strings"

	"github.com/JoshVarga/svgparser"
	"github.com/JoshVarga/svgparser/utils"
	vec "github.com/kaneta1992/vector/vector2"
)

type Glyph struct {
	Contours    []*Contour
	LeftTop     *vec.Vector2
	RightBottom *vec.Vector2
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

func (g *Glyph) updateBounds(p *vec.Vector2) {
	g.LeftTop.X = math.Min(g.LeftTop.X, p.X)
	g.LeftTop.Y = math.Min(g.LeftTop.Y, p.Y)
	g.RightBottom.X = math.Max(g.RightBottom.X, p.X)
	g.RightBottom.Y = math.Max(g.RightBottom.Y, p.Y)
}

func (g *Glyph) CreateFromSvg(filepath string) {
	svg, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal("io error")
	}
	reader := strings.NewReader(string(svg))
	element, _ := svgparser.Parse(reader, false)
	d := element.Children[0].Children[0].Attributes["d"]
	path, _ := utils.PathParser(d)
	fmt.Printf("Number of subpaths: %d\n", len(path.Subpaths))

	g.LeftTop = &vec.Vector2{999999.9, 999999.9}
	g.RightBottom = &vec.Vector2{-999999.9, -999999.9}

	// TODO: 頂点をみて自動で正規化する
	irate := 1.0 / 5000.0
	for i, subpath := range path.Subpaths {
		fmt.Printf("Path %d: ", i)
		contour := &Contour{}
		for _, command := range subpath.Commands {
			param := []float64(command.Params)
			// TODO: 構造体化する
			switch command.Symbol {
			case "M":
				v1 := &vec.Vector2{param[0], param[1]}
				v1 = v1.MulScalar(irate)
				contour.ToMove(v1)
				g.updateBounds(v1)
			case "L":
				v1 := &vec.Vector2{param[0], param[1]}
				v1 = v1.MulScalar(irate)
				contour.ToLine(v1)
				g.updateBounds(v1)
			case "Q":
				v1 := &vec.Vector2{param[0], param[1]}
				v2 := &vec.Vector2{param[2], param[3]}
				v1 = v1.MulScalar(irate)
				v2 = v2.MulScalar(irate)
				contour.ToCurve(v1, v2)
				g.updateBounds(v1)
				g.updateBounds(v2)
			default:
			}
		}
		g.AddContour(contour)
	}
}
