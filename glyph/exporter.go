package path

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/kaneta1992/vector/vector2"
	"github.com/pradeep-pyro/triangle"
)

type GlyphShader struct {
	fontPath string
	glyphs   map[rune]*Glyph
	strings  []string
}

func NewExporter(fontPath string) *GlyphShader {
	gs := &GlyphShader{fontPath: fontPath, glyphs: map[rune]*Glyph{}}
	return gs
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (g *GlyphShader) createGlyphs() {
	// 追加されてる文字列のルーン毎にGlyphを生成してmapに入れる
	for _, s := range g.strings {
		for _, r := range s {
			command := fmt.Sprintf("./bin/font.exe %s %d > tmp.svg", g.fontPath, r)
			err := exec.Command("sh", "-c", command).Run()
			if err != nil {
				fmt.Println("Command Exec Error.")
			}
			glyph := &Glyph{}
			glyph.CreateFromSvg("tmp.svg")
			g.glyphs[r] = glyph
		}
	}
}

func (g *GlyphShader) AddString(str string) {
	g.strings = append(g.strings, str)
}

func (g *GlyphShader) CreateShaderCode() string {
	g.createGlyphs()
	template_B1 := "    IB(%d,%d,%d)"
	template_B2 := "    IB2(%d,%d,%d)"
	template_T := "    IT(%d,%d,%d)"
	template_VEC := "    vec2 v[%d] = vec2[%d](\n"
	str := ""
	for _, glyph := range g.glyphs {
		points, ss, ho, beziers := glyph.CreatePointsAndInnerSegments()
		log.Print(points)
		log.Print(ss)
		log.Print(ho)
		log.Print(beziers)

		verts, faces := triangle.ConstrainedDelaunay(vec.Vec2ToFloat64(points), ss, vec.Vec2ToFloat64(ho))
		log.Print("-------------------------------")
		log.Print(verts)
		log.Print(faces)

		vertsStr := []string{}
		for _, v := range points {
			// GLSLにポートするのでYを反転する
			v.Y *= -1.0
			vertsStr = append(vertsStr, v.ToGLSLString(4))
		}

		num := len(vertsStr)

		str += fmt.Sprintf(template_VEC, num, num)
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
		for _, f := range faces {
			str += fmt.Sprintf(template_T, f[0], f[1], f[2])
			str += "\n"
		}
		for _, b := range beziers {
			v0 := points[b[0]]
			v1 := points[b[1]]
			v2 := points[b[2]]
			area := signedArea([]*vec.Vector2{v0, v1, v2})
			if area < 0.0 {
				str += fmt.Sprintf(template_B2, b[0], b[1], b[2])
			} else {
				str += fmt.Sprintf(template_B1, b[0], b[1], b[2])
			}
			str += "\n"
		}
	}

	return str
}
