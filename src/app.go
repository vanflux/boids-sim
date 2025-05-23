package boids

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	boids []Boid
}

func NewGame() *Game {
	g := &Game{
		boids: []Boid{},
	}
	for i := 0; i < 1000; i++ {
		g.boids = append(g.boids, NewBoid(rand.Float32()*500, rand.Float32()*500))
	}
	return g
}

func (g *Game) Update() error {
	for i := range g.boids {
		boid := &g.boids[i]
		boid.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	for i := range g.boids {
		boid := &g.boids[i]
		boid.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Start() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boids sim")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
