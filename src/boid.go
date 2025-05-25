package boids

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	boidImg       = ebiten.NewImage(32, 32)
)

func init() {
	whiteImage.Fill(color.White)
	initBoidAsset()
}

func initBoidAsset() {
	vertices := []ebiten.Vertex{}
	vertices = append(vertices, ebiten.Vertex{DstX: 0, DstY: 0, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 0, ColorB: 0, ColorA: 1})
	vertices = append(vertices, ebiten.Vertex{DstX: 32, DstY: 16, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 0, ColorB: 0, ColorA: 1})
	vertices = append(vertices, ebiten.Vertex{DstX: 0, DstY: 32, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 0, ColorB: 0, ColorA: 1})
	vertices = append(vertices, ebiten.Vertex{DstX: 0, DstY: 0, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 0, ColorB: 0, ColorA: 1})
	indices := []uint16{0, 1, 2, 0, 2, 3}
	boidImg.DrawTriangles(vertices, indices, whiteSubImage, &ebiten.DrawTrianglesOptions{AntiAlias: true, FillRule: ebiten.FillRuleEvenOdd})
}

type Boid struct {
	game            *Game
	x               float64
	y               float64
	vx              float64
	vy              float64
	separationRange float64
	viewRange       float64
}

func NewBoid(game *Game, x float64, y float64, vx float64, vy float64) Boid {
	boid := Boid{game: game, x: x, y: y, vx: vx, vy: vy, separationRange: 30, viewRange: 70}
	return boid
}

func (b *Boid) Update() {
	separationFactor := 0.005
	alignmentFactor := 0.05
	cohesionFactor := 0.0005
	turnFactor := 0.2
	minSpeed := 1.0
	maxSpeed := 1.5

	separationDx := 0.0
	separationDy := 0.0
	vyAvg := 0.0
	vxAvg := 0.0
	xAvg := 0.0
	yAvg := 0.0
	neighbors := 0
	for i := range b.game.boids {
		b2 := &b.game.boids[i]
		if b2 == b {
			continue
		}
		distance := math.Sqrt(((b2.x - b.x) * (b2.x - b.x)) + ((b2.y - b.y) * (b2.y - b.y))) // Pythagoras △
		if distance > b.viewRange {
			continue
		}
		neighbors++
		if distance < b.separationRange {
			separationDx += b.x - b2.x
			separationDy += b.y - b2.y
		}
		vxAvg += b2.vx
		vyAvg += b2.vy
		xAvg += b2.x
		yAvg += b2.y
	}
	if neighbors > 0 {
		vxAvg /= float64(neighbors)
		vyAvg /= float64(neighbors)
		xAvg /= float64(neighbors)
		yAvg /= float64(neighbors)

		b.vx += separationDx * separationFactor
		b.vy += separationDy * separationFactor
		b.vx += vxAvg * alignmentFactor
		b.vy += vyAvg * alignmentFactor
		b.vx += (xAvg - b.x) * cohesionFactor
		b.vy += (yAvg - b.y) * cohesionFactor
	}

	if b.x < 0 {
		b.vx += turnFactor
	}
	if b.y < 0 {
		b.vy += turnFactor
	}
	if b.x > screenWidth {
		b.vx -= turnFactor
	}
	if b.y > screenHeight {
		b.vy -= turnFactor
	}

	speed := math.Sqrt((b.vx * b.vx) + (b.vy * b.vy)) // Pythagoras △
	if speed < minSpeed {
		b.vx = b.vx / speed * minSpeed
		b.vy = b.vy / speed * minSpeed
	}
	if speed > maxSpeed {
		b.vx = b.vx / speed * maxSpeed
		b.vy = b.vy / speed * maxSpeed
	}

	b.x += b.vx
	b.y += b.vy
}

func (b *Boid) Draw(screen *ebiten.Image) {
	// Draw boid itself
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(boidImg.Bounds().Dx())/2, -float64(boidImg.Bounds().Dy())/2)
	op.GeoM.Rotate(math.Atan2(b.vy, b.vx))
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(boidImg, &op)

	// Draw angle line
	vector.StrokeLine(screen, float32(b.x), float32(b.y), float32(b.x+b.vx*20), float32(b.y+b.vy*20), 1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, false)

	// Draw view range
	// vector.StrokeCircle(screen, float32(b.x), float32(b.y), float32(b.viewRange), 1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, false)
}
