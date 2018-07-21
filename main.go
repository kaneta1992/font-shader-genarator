package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/pradeep-pyro/triangle"
)

var v float64 = 512.0

var pts = [][2]float64{}
var segs = [][2]int32{}
var holes = [][2]float64{}

var gc *draw2dimg.GraphicContext

func drawFaces(verts [][2]float64, faces [][3]int32) {
	gc.SetStrokeColor(color.RGBA{50, 50, 50, 255})
	gc.SetLineWidth(2.0)
	for _, f := range faces {
		p0 := verts[f[0]]
		p1 := verts[f[1]]
		p2 := verts[f[2]]

		x1 := p0[0]*v + v
		y1 := p0[1]*v + v
		x2 := p1[0]*v + v
		y2 := p1[1]*v + v
		x3 := p2[0]*v + v
		y3 := p2[1]*v + v

		gc.MoveTo(x1, y1)
		gc.LineTo(x2, y2)
		gc.LineTo(x3, y3)
		gc.LineTo(x1, y1)
		gc.Stroke()
	}
}

func drawSegs(verts [][2]float64, segs [][2]int32) {
	gc.SetStrokeColor(color.RGBA{50, 50, 255, 255})
	gc.SetLineWidth(2.0)
	for _, s := range segs {
		p0 := verts[s[0]]
		p1 := verts[s[1]]

		x1 := p0[0]*v + v
		y1 := p0[1]*v + v
		x2 := p1[0]*v + v
		y2 := p1[1]*v + v

		gc.MoveTo(x1, y1)
		gc.LineTo(x2, y2)
		gc.Stroke()
	}
}

func drawPts(verts [][2]float64) {
	gc.SetStrokeColor(color.RGBA{255, 50, 50, 255})
	gc.SetLineWidth(3.0)
	for _, p := range verts {
		x1 := p[0]*v + v
		y1 := p[1]*v + v

		gc.MoveTo(x1-10, y1)
		gc.LineTo(x1+10, y1)
		gc.MoveTo(x1, y1-10)
		gc.LineTo(x1, y1+10)

		gc.Stroke()
	}
}

func wrap(x, y int) int {
	return ((x % y) + y) % y
}

func signedArea(points [][2]float64) float64 {
	lp := len(points)
	area := 0.0
	for i, _ := range points {
		v1 := points[i]
		v2 := points[wrap(i+1, lp)]
		area += v1[0]*v2[1] - v1[1]*v2[0]
	}
	return area
}

func subVec2(v1, v2 [2]float64) [2]float64 {
	return [2]float64{v1[0] - v2[0], v1[1] - v2[1]}
}

func addVec2(v1, v2 [2]float64) [2]float64 {
	return [2]float64{v1[0] + v2[0], v1[1] + v2[1]}
}

func mulVec2(v [2]float64, s float64) [2]float64 {
	return [2]float64{v[0] * s, v[1] * s}
}

func length(v [2]float64) float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1])
}

func normalize(v [2]float64) [2]float64 {
	l := length(v)
	return [2]float64{v[0] / l, v[1] / l}
}

func getPath(points [][2]float64) ([][2]float64, [][2]int32, [][2]float64) {
	nextIndex := len(pts)
	lp := len(points)
	retPts := [][2]float64{}
	retSegs := [][2]int32{}
	retHoles := [][2]float64{}

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
		e0 := subVec2(v0, v1)
		e1 := subVec2(v2, v1)
		// ハーフベクトル
		hv := normalize(addVec2(e0, e1))
		// 起点頂点のハーフベクトルに少し動かした場所を穴とする
		retHoles = append(retHoles, addVec2(v1, mulVec2(hv, 0.00001)))
	}
	return retPts, retSegs, retHoles
}

func putsRect(x, y, w, h float64) {
	hw := w / 2.0
	hh := h / 2.0
	// 右回りで配置する
	p, s, ho := getPath([][2]float64{
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
	p, s, ho := getPath([][2]float64{
		{x - hw, y - hh}, {x - hw, y + hh}, {x + hw, y + hh}, {x + hw, y - hh},
	})
	pts = append(pts, p...)
	segs = append(segs, s...)
	holes = append(holes, ho...)
}

func main() {
	holes = [][2]float64{
		{99999.9, 99999.9}, // 穴がない時用の点
	}
	cutsRect(0.0, 0.0, 0.1, 0.1)
	cutsRect(0.0, 0.0, 0.025, 0.025)
	cutsRect(-0.2, -0.3, 0.2, 0.1)
	putsRect(0.0, 0.0, 0.2, 0.2)
	putsRect(0.0, 0.0, 0.05, 0.05)
	putsRect(0.2, 0.3, 0.3, 0.2)
	putsRect(-0.2, -0.3, 0.3, 0.2)

	verts, faces := triangle.ConstrainedDelaunay(pts, segs, holes)
	log.Print(pts)
	log.Print(segs)
	log.Print(holes)
	log.Print("-------------------------------")
	log.Print(verts)
	log.Print(faces)

	rgba := image.NewRGBA(image.Rect(0, 0, int(v*2.0), int(v*2.0)))

	gc = draw2dimg.NewGraphicContext(rgba)

	drawFaces(verts, faces)
	//drawSegs(verts, segs)
	drawPts(verts)
	//drawPts(holes)

	outfile, _ := os.Create("out.png")
	defer outfile.Close()
	// 書き出し
	png.Encode(outfile, rgba)
}
