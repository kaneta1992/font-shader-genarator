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

var pts = []*vec.Vector2{}
var segs = [][2]int32{}
var holes = []*vec.Vector2{}

var gc *draw2dimg.GraphicContext

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

func drawPts(verts []*vec.Vector2) {
	area := signedArea(verts)
	gc.SetStrokeColor(color.RGBA{255, 50, 50, 255})
	if area < 0.0 {
		gc.SetStrokeColor(color.RGBA{50, 255, 50, 255})
	}
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

func getPath(points []*vec.Vector2) ([]*vec.Vector2, [][2]int32, []*vec.Vector2) {
	nextIndex := len(pts)
	lp := len(points)
	retPts := []*vec.Vector2{}
	retSegs := [][2]int32{}
	retHoles := []*vec.Vector2{}

	for i, p := range points {
		retPts = append(retPts, p)
		retSegs = append(retSegs, [2]int32{int32(nextIndex + wrap(i-1, lp)), int32(nextIndex + i)})
	}

	// 左回りのパスは切り抜き用の穴を設定する
	area := signedArea(points)
	if area < 0.0 {
		// 穴を置く起点の頂点
		v0 := points[0]
		v1 := points[1]
		v2 := points[2]
		// 各頂点へのベクトル
		e0 := v0.Sub(v1)
		e1 := v2.Sub(v1)
		// ハーフベクトル
		hv := e0.Add(e1).Normalize()
		// 起点頂点のハーフベクトルに少し動かした場所を穴とする
		triArea := signedArea([]*vec.Vector2{v0, v1, v2})
		if triArea < 0.0 {
			// 起点の三角形が左回りなら内向きなのでハーフベクトル方向へ
			retHoles = append(retHoles, v1.Add(hv.MulScalar(0.00001)))
		} else {
			// 右回りなら外向きなのでハーフベクトルの逆方向へ
			retHoles = append(retHoles, v1.Sub(hv.MulScalar(0.00001)))
		}
	}
	return retPts, retSegs, retHoles
}

func putsRect(x, y, w, h float64) {
	hw := w / 2.0
	hh := h / 2.0
	// 右回りで配置する
	p, s, ho := getPath([]*vec.Vector2{
		{x - hw, y - hh}, {x + hw, y - hh}, {x + hw, y + hh}, {x - hw, y + hh},
	})
	pts = append(pts, p...)
	segs = append(segs, s...)
	holes = append(holes, ho...)
}

func cutsRect(x, y, w, h float64) {
	hw := w / 2.0
	hh := h / 2.0
	// 左回りで配置する
	p, s, ho := getPath([]*vec.Vector2{
		{x - hw, y - hh}, {x - hw, y + hh}, {x + hw, y + hh}, {x + hw, y - hh},
	})
	pts = append(pts, p...)
	segs = append(segs, s...)
	holes = append(holes, ho...)
}

func putPath(points []*vec.Vector2) {
	p, s, h := getPath(points)
	pts = append(pts, p...)
	segs = append(segs, s...)
	holes = append(holes, h...)
}

func main() {
	holes = []*vec.Vector2{
		{99999.9, 99999.9}, // 穴がない時用の点
	}

	svg, err := ioutil.ReadFile("q.svg")
	if err != nil {
		log.Fatal("io error")
	}
	reader := strings.NewReader(string(svg))
	element, _ := svgparser.Parse(reader, false)
	d := element.Children[0].Children[0].Attributes["d"]
	path, _ := utils.PathParser(d)
	fmt.Printf("Number of subpaths: %d\n", len(path.Subpaths))

	rgba := image.NewRGBA(image.Rect(0, 0, int(v*2.0), int(v*2.0)))

	gc = draw2dimg.NewGraphicContext(rgba)

	// TODO: 頂点をみて自動で正規化する
	irate := 1.0 / 5000.0
	glyph := &tpath.Glyph{}
	for i, subpath := range path.Subpaths {
		fmt.Printf("Path %d: ", i)
		points := []*vec.Vector2{}
		var nowPoint *vec.Vector2
		contour := &tpath.Contour{}
		for _, command := range subpath.Commands {
			param := []float64(command.Params)
			// TODO: 構造体化する
			switch command.Symbol {
			case "M":
				v1 := &vec.Vector2{param[0], param[1]}
				points = append(points, v1.MulScalar(irate))
				nowPoint = v1
				contour.ToMove(v1.MulScalar(irate))
			case "L":
				v1 := &vec.Vector2{param[0], param[1]}
				points = append(points, v1.MulScalar(irate))
				nowPoint = v1
				contour.ToLine(v1.MulScalar(irate))
			case "Q":
				v1 := &vec.Vector2{param[0], param[1]}
				v2 := &vec.Vector2{param[2], param[3]}
				area := signedArea([]*vec.Vector2{nowPoint, v1, v2})
				if area < 0.0 {
					// 左回りなら制御点が内部(右)にいるので制御点も含めた頂点を追加する
					points = append(points, []*vec.Vector2{v1.MulScalar(irate), v2.MulScalar(irate)}...)
				} else {
					// 右回りなら制御点が外部(左)にいるので制御点を無視する
					points = append(points, v2.MulScalar(irate))
				}
				nowPoint = v2
				contour.ToCurve(v1.MulScalar(irate), v2.MulScalar(irate))
			default:
			}
		}
		// TODO: 曲線の制御点のインデックスを考慮して三角分割する
		fmt.Println(points)
		putPath(points)
		glyph.AddContour(contour)
	}
	ps, ss := glyph.CreatePointsAndInnerSegments()
	log.Print(ps)
	log.Print(ss)
	log.Print(pts)
	log.Print(segs)
	log.Print("------------------")
	pts = ps
	segs = ss

	//cutsRect(0.0, 0.0, 0.1, 0.1)
	//cutsRect(0.0, 0.0, 0.025, 0.025)
	//cutsRect(-0.2, -0.3, 0.2, 0.1)
	//putsRect(0.0, 0.0, 0.2, 0.2)
	//putsRect(0.0, 0.0, 0.05, 0.05)
	//putsRect(0.2, 0.3, 0.3, 0.2)
	//putsRect(-0.2, -0.3, 0.3, 0.2)

	verts, faces := triangle.ConstrainedDelaunay(vec.Vec2ToFloat64(pts), segs, vec.Vec2ToFloat64(holes))
	log.Print(pts)
	log.Print(segs)
	log.Print(holes)
	log.Print("-------------------------------")
	log.Print(verts)
	log.Print(faces)

	drawFaces(vec.Float64ToVec2(verts), faces)
	drawPts(holes)
	drawPts(pts)
	drawSegs(vec.Float64ToVec2(verts), segs)

	outfile, _ := os.Create("out.png")
	defer outfile.Close()
	// 書き出し
	png.Encode(outfile, rgba)
}
