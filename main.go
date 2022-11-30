package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
	"sync"
)

const (
	screenWidth, screenHeight = 640, 360
	boidCount                 = 500
	viewRadius                = 13
	adjRate                   = 0.015
)

var (
	green   = color.RGBA{10, 255, 50, 255}
	boids   [boidCount]*Boid
	boidMap [screenWidth + 1][screenHeight + 1]int
	lock    = sync.RWMutex{}
)

type Game struct{}

func (g *Game) Update(_ *ebiten.Image) error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}
	for _, boid := range boids {
		screen.Set(int(boid.position.x+1), int(boid.position.y), green)
		screen.Set(int(boid.position.x-1), int(boid.position.y), green)
		screen.Set(int(boid.position.x), int(boid.position.y-1), green)
		screen.Set(int(boid.position.x), int(boid.position.y+1), green)
	}
}

func (g *Game) Layout(_, _ int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	for i := 0; i < boidCount; i++ {
		CreateBoid(i)
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
