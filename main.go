package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/JoshVarga/svgparser"
	"github.com/JoshVarga/svgparser/utils"
	tpath "github.com/kaneta1992/try-triangle/glyph"
	"github.com/kaneta1992/vector/vector2"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/pradeep-pyro/triangle"
)

var v float64 = 2048.0

var gc *draw2dimg.GraphicContext

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
func wrap(x, y int) int {
	return ((x % y) + y) % y
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
func drawFaces(verts []*vec.Vector2, faces [][3]int32) {

	gc.SetStrokeColor(color.RGBA{50, 50, 50, 255})
	gc.SetLineWidth(2.0)
	for _, f := range faces {
		p0 := verts[f[0]]
		p1 := verts[f[1]]
		p2 := verts[f[2]]

		p0 = p0.MulScalar(v).AddScalar(v)
		p1 = p1.MulScalar(v).AddScalar(v)
		p2 = p2.MulScalar(v).AddScalar(v)

		gc.MoveTo(p0.X, p0.Y)
		gc.LineTo(p1.X, p1.Y)
		gc.LineTo(p2.X, p2.Y)
		gc.LineTo(p0.X, p0.Y)
		gc.Stroke()
	}
}

func drawSegs(verts []*vec.Vector2, segs [][2]int32) {
	gc.SetStrokeColor(color.RGBA{50, 50, 255, 255})
	gc.SetLineWidth(2.0)
	for _, s := range segs {
		p0 := verts[s[0]]
		p1 := verts[s[1]]

		p0 = p0.MulScalar(v).AddScalar(v)
		p1 = p1.MulScalar(v).AddScalar(v)

		gc.MoveTo(p0.X, p0.Y)
		gc.LineTo(p1.X, p1.Y)
		gc.Stroke()
	}
}

func drawPts(verts []*vec.Vector2, col color.RGBA) {
	gc.SetStrokeColor(col)
	gc.SetLineWidth(3.0)
	for _, p := range verts {
		p0 := p.MulScalar(v).AddScalar(v)

		gc.MoveTo(p0.X-10, p0.Y)
		gc.LineTo(p0.X+10, p0.Y)
		gc.MoveTo(p0.X, p0.Y-10)
		gc.LineTo(p0.X, p0.Y+10)

		gc.Stroke()
	}
}

func main() {
	svg, err := ioutil.ReadFile("a.svg")
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
	glyph := &tpath.Glyph{}
	for i, subpath := range path.Subpaths {
		fmt.Printf("Path %d: ", i)
		contour := &tpath.Contour{}
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
		glyph.AddContour(contour)
	}
	ps, ss, ho, be := glyph.CreatePointsAndInnerSegments()
	log.Print(ps)
	log.Print(ss)
	log.Print(ho)
	log.Print(be)

	verts, faces := triangle.ConstrainedDelaunay(vec.Vec2ToFloat64(ps), ss, vec.Vec2ToFloat64(ho))
	log.Print("-------------------------------")
	log.Print(verts)
	log.Print(faces)

	rgba := image.NewRGBA(image.Rect(0, 0, int(v*2.0), int(v*2.0)))
	gc = draw2dimg.NewGraphicContext(rgba)

	drawFaces(vec.Float64ToVec2(verts), faces)
	drawPts(ho, color.RGBA{50, 255, 50, 255})
	drawPts(ps, color.RGBA{255, 50, 50, 255})
	drawSegs(vec.Float64ToVec2(verts), ss)

	outfile, _ := os.Create("out.png")
	defer outfile.Close()
	// 書き出し
	png.Encode(outfile, rgba)

	template_B1 := "    d = min( d, InBezier(v[%d],v[%d],v[%d], uv));"
	template_B2 := "    d = min( d, InBezier2(v[%d],v[%d],v[%d], uv));"
	template_T := "    d = min( d, InTri(v[%d],v[%d],v[%d], uv));"
	template_VEC := "    vec2 v[%d] = vec2[%d](\n"
	//vec2 v[2] = vec2[2](
	//   	vec2(0.1244,-0.4704),vec2(0.1248,-0.5408)
	//);

	vertsStr := []string{}
	for _, v := range ps {
		vertsStr = append(vertsStr, v.ToGLSLString(4))
	}

	num := len(vertsStr)

	str := fmt.Sprintf(template_VEC, num, num)
	for i := 0; i < num/6+1; i++ {
		arr := vertsStr[i*6 : Min(i*6+6, num)]
		str += "        "
		str += strings.Join(arr, ",")
		if i != num/6 {
			str += ","
		}
		str += "\n"
	}
	str += "    );\n"
	fmt.Print(str)
	str = ""
	for _, f := range faces {
		str += fmt.Sprintf(template_T, f[2], f[1], f[0])
		str += "\n"
	}
	for _, b := range be {
		v0 := ps[b[0]]
		v1 := ps[b[1]]
		v2 := ps[b[2]]
		area := signedArea([]*vec.Vector2{v0, v1, v2})
		if area < 0.0 {
			str += fmt.Sprintf(template_B1, b[2], b[1], b[0])
		} else {
			str += fmt.Sprintf(template_B2, b[2], b[1], b[0])
		}
		str += "\n"
	}
	fmt.Print(str)
}
