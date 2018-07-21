package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/pradeep-pyro/triangle"
)

var v float64 = 512.0
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

func main() {
	var pts = [][2]float64{{-0.1, -0.1}, {-0.5, -0.1}, {-0.5, -0.5}, {-0.1, -0.5}, {0.1, 0.1}, {0.5, 0.1}, {0.5, 0.5}, {0.1, 0.5},
		{0.15, 0.15}, {0.45, 0.15}, {0.45, 0.45}, {0.15, 0.45}, {0.2, 0.2}, {0.4, 0.2}, {0.4, 0.4}, {0.2, 0.4},
		{0.25, 0.25}, {0.35, 0.25}, {0.35, 0.35}, {0.25, 0.35},
	}
	var segs = [][2]int32{{3, 0}, {0, 1}, {1, 2}, {2, 3}, {7, 4}, {4, 5}, {5, 6}, {6, 7}, {11, 8}, {8, 9}, {9, 10}, {10, 11}, {15, 12}, {12, 13}, {13, 14}, {14, 15}, {19, 16}, {16, 17}, {17, 18}, {18, 19}}
	var holes = [][2]float64{
		{99999.9, 99999.9}, // 穴がない時用の点
		{0.16, 0.44},
		{0.26, 0.34},
	}
	verts, faces := triangle.ConstrainedDelaunay(pts, segs, holes)
	log.Print(pts)
	log.Print(segs)
	log.Print(verts)
	log.Print(faces)

	rgba := image.NewRGBA(image.Rect(0, 0, int(v*2.0), int(v*2.0)))

	gc = draw2dimg.NewGraphicContext(rgba)

	drawFaces(verts, faces)
	//drawSegs(verts, segs)
	drawPts(verts)

	outfile, _ := os.Create("out.png")
	defer outfile.Close()
	// 書き出し
	png.Encode(outfile, rgba)
}
