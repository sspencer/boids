package main

import (
	"boid/single"
	"boid/threaded"
	"boid/util"
	"flag"
	"fmt"
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

var (
	col = color.RGBA{R: 10, G: 255, B: 50, A: 255}
)

type World interface {
	Setup(width, height, count, radius int, rate float64)
	Animate()
	PositionAndVelocity(id int) (util.Vector2D, util.Vector2D)
}

type Game struct {
	world World
	count int
}

func (g *Game) Update(_ *ebiten.Image) error {
	g.world.Animate()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < g.count; i++ {
		pos, _ := g.world.PositionAndVelocity(i)
		screen.Set(int(pos.X+1), int(pos.Y), col)
		screen.Set(int(pos.X-1), int(pos.Y), col)
		screen.Set(int(pos.X), int(pos.Y-1), col)
		screen.Set(int(pos.X), int(pos.Y+1), col)
	}
}

func (g *Game) Layout(_, _ int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	var threadedWorld bool
	var count, radius int
	var adjustment float64

	flag.BoolVar(&threadedWorld, "t", false, "One thread per boid")
	flag.IntVar(&count, "n", boidCount, "Number of boids")
	flag.IntVar(&radius, "r", viewRadius, "View radius")
	flag.Float64Var(&adjustment, "a", adjRate, "Adjustment rate, smaller is smoother")
	flag.Parse()

	var world World

	if threadedWorld {
		fmt.Println("Started Multi Threaded World")
		world = threaded.NewBoidWorld()
	} else {
		fmt.Println("Started Single Threaded World")
		world = single.NewBoidWorld()
	}

	world.Setup(screenWidth, screenHeight, count, radius, adjustment)
	game := Game{world: world, count: count}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("World")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
