package path

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/JoshVarga/svgparser"
	"github.com/JoshVarga/svgparser/utils"
	vec "github.com/kaneta1992/vector/vector2"
)

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
				contour.ToMove(v1.MulScalar(irate))
			case "L":
				v1 := &vec.Vector2{param[0], param[1]}
				contour.ToLine(v1.MulScalar(irate))
			case "Q":
				v1 := &vec.Vector2{param[0], param[1]}
				v2 := &vec.Vector2{param[2], param[3]}
				contour.ToCurve(v1.MulScalar(irate), v2.MulScalar(irate))
			default:
			}
		}
		g.AddContour(contour)
	}
}
