package main

import (
	"boid/threads"
	"boid/util"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 360
	boidCount    = 500
	viewRadius   = 13
	adjRate      = 0.015
)

type World interface {
	Setup(width, height, count, radius int, rate float64)
	Animate()
	Position(id int) util.Vector2D
}

var (
	green = color.RGBA{R: 10, G: 255, B: 50, A: 255}
)

type Game struct {
	world World
	count int
}

func (g *Game) Update(_ *ebiten.Image) error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Animate()
	for i := 0; i < g.count; i++ {
		pos := g.world.Position(i)
		screen.Set(int(pos.X+1), int(pos.Y), green)
		screen.Set(int(pos.X-1), int(pos.Y), green)
		screen.Set(int(pos.X), int(pos.Y-1), green)
		screen.Set(int(pos.X), int(pos.Y+1), green)
	}
}

func (g *Game) Layout(_, _ int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	world := threads.NewBoidWorld()
	world.Setup(screenWidth, screenHeight, boidCount, viewRadius, adjRate)
	game := Game{world: world, count: boidCount}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("World")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
