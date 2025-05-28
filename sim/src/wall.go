package boids

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Wall struct {
	game   *Game
	x      float64
	y      float64
	width  float64
	height float64
}

func NewWall(game *Game, x float64, y float64, width float64, height float64) Wall {
	wall := Wall{game: game, x: x, y: y, width: width, height: height}
	return wall
}

func (w *Wall) Update() {}

func (w *Wall) Draw(screen *ebiten.Image) {
	colorR := uint8(255)
	colorG := uint8(255)
	colorB := uint8(255)

	vector.StrokeRect(screen, float32(w.x), float32(w.y), float32(w.width), float32(w.height), 1, color.RGBA{R: colorR, G: colorG, B: colorB, A: 255}, true)
}
