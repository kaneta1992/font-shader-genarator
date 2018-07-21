package main

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/pradeep-pyro/triangle"
)

var v float64 = 500.0
var gc *draw2dimg.GraphicContext

func drawFaces(verts [][2]float64, faces [][3]int32) {
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
	for _, p := range verts {
		x1 := p[0]*v + v
		y1 := p[1]*v + v

		gc.MoveTo(x1, y1)
		gc.LineTo(x1, y1-20.0)
		gc.Stroke()
	}
}

func main() {
	// Points forming the shape of letter "A"
	var pts = [][2]float64{{0.200000, -0.776400}, {0.220000, -0.773200},
		{0.245600, -0.756400}, {0.277600, -0.702000}, {0.488800, -0.207600}, {0.504800, -0.207600}, {0.740800, -0.7396}, {0.756000, -0.761200},
		{0.774400, -0.7724}, {0.800000, -0.776400}, {0.800000, -0.792400}, {0.579200, -0.792400}, {0.579200, -0.776400}, {0.621600, -0.771600},
		{0.633600, -0.762800}, {0.639200, -0.744400}, {0.620800, -0.684400}, {0.587200, -0.604400}, {0.360800, -0.604400}, {0.319200, -0.706800},
		{0.312000, -0.739600}, {0.318400, -0.761200}, {0.334400, -0.771600}, {0.371200, -0.776400}, {0.371200, -0.792400}, {0.374400, -0.570000},
		{0.574400, -0.5700}, {0.473600, -0.330800}, {0.200000, -0.792400},
	}
	// Segments connecting the points
	var segs = [][2]int32{{28, 0}, {0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}, {12, 13}, {13, 14}, {14, 15}, {15, 16}, {16, 17}, {17, 18}, {18, 19}, {19, 20}, {20, 21}, {21, 22}, {22, 23}, {23, 24}, {24, 28}, {25, 26}, {26, 27}, {27, 25}}
	// Hole represented by a point lying inside it
	var holes = [][2]float64{
		{0.47, -0.5},
	}
	verts, faces := triangle.ConstrainedDelaunay(pts, segs, holes)
	log.Print(pts)
	log.Print(segs)
	log.Print(verts)
	log.Print(faces)

	rgba := image.NewRGBA(image.Rect(0, 0, 1024, 1024))

	gc = draw2dimg.NewGraphicContext(rgba)
	//drawPts(verts)
	//drawSegs(verts, segs)

	drawFaces(verts, faces)

	outfile, _ := os.Create("out.png")
	defer outfile.Close()
	// 書き出し
	png.Encode(outfile, rgba)
}
