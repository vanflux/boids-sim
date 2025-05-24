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
	game         *Game
	x            float64
	y            float64
	vx           float64
	vy           float64
	separationX  float64
	separationY  float64
	alignmentX   float64
	alignmentY   float64
	cohesionX    float64
	cohesionY    float64
	desiredAngle float64
	angle        float64
	viewRange    float64
}

func NewBoid(game *Game, x float64, y float64, angle float64) Boid {
	boid := Boid{game: game, x: x, y: y, separationX: 0, separationY: 0, alignmentX: 0, alignmentY: 0, cohesionX: 0, cohesionY: 0, desiredAngle: 0, angle: angle, viewRange: 100}
	return boid
}

func (b *Boid) Update() {
	b.x = b.x + b.vx
	b.y = b.y + b.vy

	b.separationX = 0
	b.separationY = 0
	b.alignmentX = 0
	b.alignmentY = 0
	b.cohesionX = 0
	b.cohesionY = 0
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
		b.separationX -= b2.x - b.x
		b.separationY -= b2.y - b.y
		b.alignmentX += b2.vx
		b.alignmentY += b2.vy
		b.cohesionX += b2.x - b.x
		b.cohesionY += b2.y - b.y
	}
	separationAngle := math.Atan2(b.separationY, b.separationX)
	b.separationX = math.Cos(separationAngle)
	b.separationY = math.Sin(separationAngle)
	alignmentAngle := math.Atan2(b.alignmentY, b.alignmentX)
	b.alignmentX = math.Cos(alignmentAngle)
	b.alignmentY = math.Sin(alignmentAngle)
	centroidX := b.cohesionX / float64(neighbors)
	centroidY := b.cohesionY / float64(neighbors)
	cohesionAngle := math.Atan2(centroidY, centroidX)
	b.cohesionX = math.Cos(cohesionAngle)
	b.cohesionY = math.Sin(cohesionAngle)

	if neighbors > 0 {
		// b.desiredAngle = math.Atan2(b.cohesionY, b.cohesionX)
		print(b.separationY+b.alignmentY+b.cohesionY, b.separationX+b.alignmentX+b.cohesionX, "\n")
		b.desiredAngle = math.Atan2(b.separationY+b.alignmentY+b.cohesionY, b.separationX+b.alignmentX+b.cohesionX)
		shortestAngle := math.Mod((b.desiredAngle-b.angle)+math.Pi, math.Pi*2) - math.Pi
		stepAngle := math.Min(math.Max(-0.05, shortestAngle), 0.05)
		b.angle += stepAngle
	}

	// b.angle = math.Mod((b.angle + (rand.Float64()-0.5)*0.5), math.Pi*2)
	b.vx = math.Cos(b.angle)
	b.vy = math.Sin(b.angle)
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

	// Draw separation line
	vector.StrokeLine(screen, float32(b.x), float32(b.y), float32(b.x+b.separationX*20), float32(b.y+b.separationY*20), 1, color.RGBA{R: 255, G: 255, B: 0, A: 255}, false)

	// Draw aligment line
	vector.StrokeLine(screen, float32(b.x), float32(b.y), float32(b.x+b.alignmentX*20), float32(b.y+b.alignmentY*20), 1, color.RGBA{R: 0, G: 255, B: 0, A: 255}, false)

	// Draw cohesion line
	vector.StrokeLine(screen, float32(b.x), float32(b.y), float32(b.x+b.cohesionX*20), float32(b.y+b.cohesionY*20), 1, color.RGBA{R: 0, G: 255, B: 255, A: 255}, false)

	// Draw view range
	vector.StrokeCircle(screen, float32(b.x), float32(b.y), float32(b.viewRange), 1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, false)
}
