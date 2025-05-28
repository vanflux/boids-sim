package boids

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1200
	screenHeight = 700
	normalBoids  = 100
	enemyBoids   = 3
	walls        = 3
	sampleRate   = 44100
)

type Game struct {
	boids []Boid
	walls []Wall
}

func NewGame() (*Game, error) {
	g := &Game{
		boids: []Boid{},
	}
	for i := 0; i < normalBoids; i++ {
		boid, _ := NewBoid(g, Normal, rand.Float64()*screenWidth, rand.Float64()*screenHeight, rand.Float64()-0.5, rand.Float64()-0.5)
		g.boids = append(g.boids, *boid)
	}
	for i := 0; i < enemyBoids; i++ {
		boid, _ := NewBoid(g, Enemy, rand.Float64()*screenWidth, rand.Float64()*screenHeight, rand.Float64()-0.5, rand.Float64()-0.5)
		g.boids = append(g.boids, *boid)
	}
	// for i := 0; i < walls; i++ {
	// 	g.walls = append(g.walls, NewWall(g, rand.Float64()*(screenWidth-80)+40, rand.Float64()*(screenHeight-80)+40, rand.Float64()*20+20, rand.Float64()*20+20))
	// }
	g.walls = append(g.walls, NewWall(g, 40, 40, 100, 20))
	g.walls = append(g.walls, NewWall(g, 300, 50, 80, 80))
	g.walls = append(g.walls, NewWall(g, 200, 300, 20, 150))
	return g, nil
}

func (g *Game) detectWall(x float64, y float64) bool {
	for i := range g.walls {
		w := &g.walls[i]
		if x > w.x && y > w.y && x < w.x+w.width && y < w.y+w.height {
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
	for i := range g.boids {
		boid := &g.boids[i]
		boid.Update()
	}
	for i := range g.walls {
		wall := &g.walls[i]
		wall.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	for i := range g.boids {
		boid := &g.boids[i]
		boid.Draw(screen)
	}
	for i := range g.walls {
		wall := &g.walls[i]
		wall.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Start() error {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boids sim")
	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
		return err
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
