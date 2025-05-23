package boids

import (
	"image"
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

type Boid struct {
	x    float32
	y    float32
	vx   float32
	vy   float32
	size float32
}

func NewBoid(x float32, y float32) Boid {
	boid := Boid{x: x, y: y, size: 10}
	return boid
}

func (b *Boid) Update() {
	b.x = b.x + b.vx
	b.y = b.y + b.vy
	b.vx = rand.Float32()*4 - 2
	b.vy = rand.Float32()*4 - 2
}

func (b *Boid) Draw(screen *ebiten.Image) {
	var path vector.Path

	// E
	path.MoveTo(0, 0)
	path.LineTo(50, 25)
	path.LineTo(0, 50)
	path.LineTo(0, 0)

	vertices := []ebiten.Vertex{}
	indices := []uint16{}

	vertices, indices = path.AppendVerticesAndIndicesForStroke(vertices[:0], indices[:0], &vector.StrokeOptions{
		Width: 1,
	})

	for i := range vertices {
		vertices[i].DstX = (vertices[i].DstX + float32(0))
		vertices[i].DstY = (vertices[i].DstY + float32(0))
		vertices[i].SrcX = 1
		vertices[i].SrcY = 1
		vertices[i].ColorR = 0xdb / float32(0xff)
		vertices[i].ColorG = 0x56 / float32(0xff)
		vertices[i].ColorB = 0x20 / float32(0xff)
		vertices[i].ColorA = 1
		// print("vertices[i]:", i, "    ", vertices[i].DstX, ", ", vertices[i].DstY, ", ", vertices[i].SrcX, ", ", vertices[i].SrcY, ", ", vertices[i].ColorR, ", ", vertices[i].ColorG, ", ", vertices[i].ColorB, ", ", vertices[i].ColorA, ", ", vertices[i].DstX, ", ", vertices[i].DstX, "    \n\n")
	}

	objImg := ebiten.NewImage(80, 80)
	objImg.DrawTriangles(vertices, indices, whiteSubImage, &ebiten.DrawTrianglesOptions{AntiAlias: true, FillRule: ebiten.FillRuleEvenOdd})

	op := ebiten.DrawImageOptions{}
	// op.GeoM.Rotate(math.Pi / 180 * 45)
	// op.GeoM.Translate(100, 100)
	screen.DrawImage(objImg, &op)

	// vector.StrokeLine(screen, b.x, b.y, b.x+b.size, b.y+b.size/2, 1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, true)
	// vector.StrokeLine(screen, b.x, b.y, b.x, b.y+b.size, 1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, true)
	// vector.StrokeLine(screen, b.x, b.y+b.size, b.x+b.size, b.y+b.size/2, 1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, true)
}
