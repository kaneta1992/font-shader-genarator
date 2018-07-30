package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	tpath "github.com/kaneta1992/try-triangle/glyph"
	"github.com/kaneta1992/vector/vector2"
	"github.com/llgcode/draw2d/draw2dimg"
)

var v float64 = 2048.0

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
	tpath.CW = false
	gs := tpath.NewExporter("ns.ttf")
	gs.AddString("欺く為のフェイントです")

	rgba := image.NewRGBA(image.Rect(0, 0, int(v*2.0), int(v*2.0)))
	gc = draw2dimg.NewGraphicContext(rgba)

	// drawFaces(vec.Float64ToVec2(verts), faces)
	// drawPts(ho, color.RGBA{50, 255, 50, 255})
	// drawPts(ps, color.RGBA{255, 50, 50, 255})
	// drawSegs(vec.Float64ToVec2(verts), ss)

	outfile, _ := os.Create("out.png")
	defer outfile.Close()
	// 書き出し
	png.Encode(outfile, rgba)

	fmt.Print(gs.CreateShaderCode())
}
