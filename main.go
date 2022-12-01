package main

import (
	"boid/single"
	"boid/threaded"
	"boid/util"
	"flag"
	"fmt"
	gc "github.com/gerow/go-color"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
)

const (
	screenWidth     = 640
	screenHeight    = 360
	boidCount       = 800
	viewRadius      = 13
	singleAdjRate   = 0.25
	threadedAdjRate = 0.015
)

var (
	col = color.RGBA{R: 10, G: 255, B: 50, A: 255}
)

type World interface {
	Setup(width, height, count, radius int, rate float64)
	Animate()
	Position(id int) *util.Vector2D
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
	mid := util.Vector2D{X: screenWidth / 2, Y: screenHeight / 2}
	maxDist := mid.Distance(util.Vector2D{X: 0, Y: 0})

	for i := 0; i < g.count; i++ {
		pos := g.world.Position(i)

		dist := mid.Distance(*pos)

		hsl := gc.HSL{H: dist / maxDist, S: 1.0, L: 0.5}
		rgb := hsl.ToRGB()
		col = color.RGBA{
			R: uint8(rgb.R * 255),
			G: uint8(rgb.G * 255),
			B: uint8(rgb.B * 255),
			A: 255,
		}

		screen.Set(int(pos.X+1), int(pos.Y), col)
		screen.Set(int(pos.X-1), int(pos.Y), col)
		screen.Set(int(pos.X), int(pos.Y), col)
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
	var adjRate float64

	flag.BoolVar(&threadedWorld, "t", false, "One thread per boid")
	flag.IntVar(&count, "n", boidCount, "Number of boids")
	flag.IntVar(&radius, "r", viewRadius, "View radius")
	flag.Parse()

	var world World

	if threadedWorld {
		adjRate = threadedAdjRate
		world = threaded.NewBoidWorld()
		fmt.Println("Started Multi Threaded World")
	} else {
		adjRate = singleAdjRate
		world = single.NewBoidWorld()
		fmt.Println("Started Single Threaded World")
	}

	world.Setup(screenWidth, screenHeight, count, radius, adjRate)
	game := Game{world: world, count: count}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
