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
func Chunks(strs []string, n int) [][]string {
	num := len(strs)
	ret := [][]string{}
	for i := 0; i <= (num-1)/n; i++ {
		ret = append(ret, strs[i*n:Min(i*n+n, num)])
	}
	return ret
}
func JoinChunks(chunks [][]string, chunkSep, strSep string) string {
	strs := []string{}
	for _, c := range chunks {
		strs = append(strs, strings.Join(c, strSep))
	}
	return strings.Join(strs, chunkSep)
}

func (g *GlyphShader) CreateShaderCode() string {
	g.createGlyphs()
	templateFunc := "float _%X(vec2 uv) {    // %s\n    float d = 10000.0;\n"
	templateB1 := "IB(%d,%d,%d)"
	templateB2 := "IB2(%d,%d,%d)"
	templateT := "IT(%d,%d,%d)"
	templateVEC := "    vec2 v[%d] = vec2[%d](\n"
	str := ""
	for r, glyph := range g.glyphs {
		points, ss, ho, beziers := glyph.CreatePointsAndInnerSegments()
		log.Print(points)
		log.Print(ss)
		log.Print(ho)
		log.Print(beziers)

		verts, faces := triangle.ConstrainedDelaunay(vec.Vec2ToFloat64(points), ss, vec.Vec2ToFloat64(ho))
		log.Print("-------------------------------")
		log.Print(verts)
		log.Print(faces)

		str += fmt.Sprintf(templateFunc, r, string(r))

		vertStrs := []string{}
		for _, v := range points {
			// GLSLにポートするのでYを反転する
			v.Y *= -1.0
			vertStrs = append(vertStrs, v.ToGLSLString(4))
		}

		num := len(vertStrs)

		str += fmt.Sprintf(templateVEC, num, num)

		vertChunks := Chunks(vertStrs, 6)
		str += "        " + JoinChunks(vertChunks, ",\n        ", ",") + "\n    );\n"

		geomStrs := []string{}
		for _, f := range faces {
			geomStrs = append(geomStrs, fmt.Sprintf(templateT, f[0], f[1], f[2]))
		}
		for _, b := range beziers {
			v0 := points[b[0]]
			v1 := points[b[1]]
			v2 := points[b[2]]
			area := signedArea([]*vec.Vector2{v0, v1, v2})
			if area < 0.0 {
				geomStrs = append(geomStrs, fmt.Sprintf(templateB2, b[0], b[1], b[2]))
			} else {
				geomStrs = append(geomStrs, fmt.Sprintf(templateB1, b[0], b[1], b[2]))
			}
		}

		geomChunks := Chunks(geomStrs, 10)
		str += "    " + JoinChunks(geomChunks, "\n    ", "") + "\n"
		str += "    return d;\n}\n"
	}

	return str
}
