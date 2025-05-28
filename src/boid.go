package boids

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	boidImg       = ebiten.NewImage(32, 32)
	audioContext  *audio.Context
	audioPlayer   *audio.Player
	lastPlay      = int64(0)
)

func init() {
	whiteImage.Fill(color.White)
	initBoidAsset()
	audioContext = audio.NewContext(44100)
	dat, err := os.ReadFile("/home/lucas/projects/boids-sim/src/hit.wav")
	if err == nil {
		s, err := wav.DecodeF32(bytes.NewReader(dat))
		if err == nil {
			audioPlayer, _ = audioContext.NewPlayerF32(s)
			audioPlayer.SetVolume(0.1)
		}
	}
}

func initBoidAsset() {
	vertices := []ebiten.Vertex{}
	vertices = append(vertices, ebiten.Vertex{DstX: 0, DstY: 0, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1})
	vertices = append(vertices, ebiten.Vertex{DstX: 32, DstY: 16, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1})
	vertices = append(vertices, ebiten.Vertex{DstX: 0, DstY: 32, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1})
	vertices = append(vertices, ebiten.Vertex{DstX: 0, DstY: 0, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1})
	indices := []uint16{0, 1, 2, 0, 2, 3}
	boidImg.DrawTriangles(vertices, indices, whiteSubImage, &ebiten.DrawTrianglesOptions{AntiAlias: true, FillRule: ebiten.FillRuleEvenOdd})
}

type BoidKind int8

const (
	Normal BoidKind = iota
	Enemy
)

type Boid struct {
	game            *Game
	kind            BoidKind
	x               float64
	y               float64
	vx              float64
	vy              float64
	separationRange float64
	viewRange       float64
}

func NewBoid(game *Game, kind BoidKind, x float64, y float64, vx float64, vy float64) (*Boid, error) {
	boid := &Boid{game: game, kind: kind, x: x, y: y, vx: vx, vy: vy, separationRange: 30, viewRange: 70}
	return boid, nil
}

func (b *Boid) raycastWall(angleDiff float64, distance float64) bool {
	forwardAngle := math.Atan2(b.vy, b.vx)
	angle := forwardAngle + angleDiff
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	for i := float64(0); i <= distance; i += 5 {
		x := b.x + cos*distance
		y := b.y + sin*distance
		if b.game.detectWall(x, y) {
			return true
		}
	}
	return false
}

func (b *Boid) applyRules() {
	separationFactor := 0.03
	alignmentFactor := 0.02
	cohesionFactor := 0.001

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
			separationAngle := math.Atan2(b.y-b2.y, b.x-b2.x)
			separationDx += math.Cos(separationAngle) * (b.separationRange - distance)
			separationDy += math.Sin(separationAngle) * (b.separationRange - distance)
		}
		if distance < b.viewRange && b2.kind == Enemy {
			separationAngle := math.Atan2(b.y-b2.y, b.x-b2.x)
			separationDx += math.Cos(separationAngle) * (b.viewRange - distance)
			separationDy += math.Sin(separationAngle) * (b.viewRange - distance)
		}
		vxAvg += b2.vx
		vyAvg += b2.vy
		xAvg += b2.x
		yAvg += b2.y
	}
	avoidWallDistance := 70.0
	avoidWallFactor := 25.0
	avoidWallMinAngle := math.Pi * 2 * 0.05
	avoidWallMaxAngle := math.Pi * 2 * 0.30
	avoidWallAngleStep := math.Pi * 2 * 0.05
loop:
	for distance := 1.0; distance <= avoidWallDistance; distance++ {
		for angleDiff := avoidWallMinAngle; angleDiff <= avoidWallMaxAngle; angleDiff += avoidWallAngleStep {
			if b.raycastWall(angleDiff, distance) {
				separationAngle := math.Atan2(b.vy, b.vx) - (math.Pi / 2)
				separationDx += math.Cos(separationAngle) * (avoidWallDistance - distance) / avoidWallDistance * avoidWallFactor
				separationDy += math.Sin(separationAngle) * (avoidWallDistance - distance) / avoidWallDistance * avoidWallFactor
				break loop
			} else if b.raycastWall(-angleDiff, distance) {
				separationAngle := math.Atan2(b.vy, b.vx) + (math.Pi / 2)
				separationDx += math.Cos(separationAngle) * (avoidWallDistance - distance) / avoidWallDistance * avoidWallFactor
				separationDy += math.Sin(separationAngle) * (avoidWallDistance - distance) / avoidWallDistance * avoidWallFactor
				break loop
			}
		}
	}
	if ebiten.IsMouseButtonPressed(0) {
		cursorX, cursorY := ebiten.CursorPosition()
		distance := math.Sqrt(((float64(cursorX) - b.x) * (float64(cursorX) - b.x)) + ((float64(cursorY) - b.y) * (float64(cursorY) - b.y))) // Pythagoras △
		if distance < b.viewRange {
			separationAngle := math.Atan2(b.y-float64(cursorY), b.x-float64(cursorX))
			separationDx += math.Cos(separationAngle) * (b.viewRange - distance)
			separationDy += math.Sin(separationAngle) * (b.viewRange - distance)
		}
	}
	b.vx += separationDx * separationFactor
	b.vy += separationDy * separationFactor
	if neighbors > 0 {
		vxAvg /= float64(neighbors)
		vyAvg /= float64(neighbors)
		xAvg /= float64(neighbors)
		yAvg /= float64(neighbors)
		b.vx += vxAvg * alignmentFactor
		b.vy += vyAvg * alignmentFactor
		b.vx += (xAvg - b.x) * cohesionFactor
		b.vy += (yAvg - b.y) * cohesionFactor
	}
}

func (b *Boid) applyInfiniteScreen() {
	tolerance := 10.0
	if b.x < -tolerance {
		b.x = screenWidth
	}
	if b.y < -tolerance {
		b.y = screenHeight
	}
	if b.x > screenWidth+tolerance {
		b.x = 0
	}
	if b.y > screenHeight+tolerance {
		b.y = 0
	}
}

func (b *Boid) applyWallsCollision() {
	for i := range b.game.walls {
		w := &b.game.walls[i]

		if b.x > w.x && b.y > w.y && b.x < w.x+w.width && b.y < w.y+w.height {
			lDist := b.x - w.x
			rDist := w.x + w.width - b.x
			tDist := b.y - w.y
			bDist := w.y + w.height - b.y

			minDist := min(lDist, rDist, tDist, bDist)
			if time.Now().UnixMilli()-lastPlay > 50 {
				lastPlay = time.Now().UnixMilli()
				audioPlayer.Rewind()
				audioPlayer.Play()
			}
			if lDist == minDist {
				b.x += w.x - b.x
				b.vx = -b.vx
			} else if rDist == minDist {
				b.x += w.x + w.width - b.x
				b.vx = -b.vx
			} else if tDist == minDist {
				b.y += w.y - b.y
				b.vy = -b.vy
			} else if bDist == minDist {
				b.y += w.y + w.height - b.y
				b.vy = -b.vy
			}
		}
	}
}

func (b *Boid) applySpeedLimit() {
	minSpeed := 1.0
	maxSpeed := 1.5
	speed := math.Sqrt((b.vx * b.vx) + (b.vy * b.vy)) // Pythagoras △
	if speed < minSpeed {
		b.vx = b.vx / speed * minSpeed
		b.vy = b.vy / speed * minSpeed
	}
	if speed > maxSpeed {
		b.vx = b.vx / speed * maxSpeed
		b.vy = b.vy / speed * maxSpeed
	}
}

func (b *Boid) Update() {
	b.applyRules()
	b.applyWallsCollision()
	b.applySpeedLimit()
	b.x += b.vx
	b.y += b.vy
	b.applyInfiniteScreen()
}

func (b *Boid) Draw(screen *ebiten.Image) {
	colorR := uint8(255)
	colorG := uint8(255)
	colorB := uint8(255)
	if b.kind == Enemy {
		colorG = 0
		colorB = 0
	}

	// Draw boid itself
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(boidImg.Bounds().Dx())/2, -float64(boidImg.Bounds().Dy())/2)
	op.GeoM.Rotate(math.Atan2(b.vy, b.vx))
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(b.x, b.y)
	op.ColorScale.SetR(float32(colorR) / 255)
	op.ColorScale.SetG(float32(colorG) / 255)
	op.ColorScale.SetB(float32(colorB) / 255)
	screen.DrawImage(boidImg, &op)

	// Draw angle line
	// vector.StrokeLine(screen, float32(b.x), float32(b.y), float32(b.x+b.vx*20), float32(b.y+b.vy*20), 1, color.RGBA{R: colorR, G: colorG, B: colorB, A: 255}, false)

	// Draw view range
	// vector.StrokeCircle(screen, float32(b.x), float32(b.y), float32(b.viewRange), 1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, false)
}
