package single

import (
	"boid/util"
	"math"
	"math/rand"
)

var (
	screenWidth  = 0.0
	screenHeight = 0.0
	viewRadius   = 0.0
	adjRate      = 0.0
	boids        []*boid
	boidMap      [][]int
)

type BoidWorld struct {
}

type boid struct {
	id       int
	position util.Vector2D
	velocity util.Vector2D
	accel    util.Vector2D
}

func NewBoidWorld() *BoidWorld {
	return &BoidWorld{}
}

func (w *BoidWorld) Setup(width, height, count, radius int, rate float64) {
	screenWidth = float64(width)
	screenHeight = float64(height)
	viewRadius = float64(radius)
	adjRate = rate

	boids = make([]*boid, count)
	boidMap = make([][]int, width+1)
	for i := range boidMap {
		boidMap[i] = make([]int, height+1)
		for j := 0; j < len(boidMap[i]); j++ {
			boidMap[i][j] = -1
		}
	}

	for i := 0; i < count; i++ {
		bid := i
		b := &boid{id: bid}
		b.calcPosition()
		b.calcVelocity()
		boids[bid] = b
		boidMap[int(b.position.X)][int(b.position.Y)] = bid
	}
}

func (w *BoidWorld) Animate() {

	for _, b := range boids {
		b.calcAcceleration()
	}

	for _, b := range boids {
		b.velocity = b.velocity.Add(b.accel).Limit(-1, 1)
		boidMap[int(b.position.X)][int(b.position.Y)] = -1
		b.position = b.position.Add(b.velocity)
		boidMap[int(b.position.X)][int(b.position.Y)] = b.id
	}
}

func (w *BoidWorld) PositionVelocity(id int) (*util.Vector2D, *util.Vector2D) {
	return &boids[id].position, &boids[id].velocity
}

func (b *boid) calcPosition() {
	b.position.X = rand.Float64() * screenWidth
	b.position.Y = rand.Float64() * screenHeight
}

func (b *boid) calcVelocity() {
	b.velocity.X = rand.Float64()*2 - 1.0
	b.velocity.Y = rand.Float64()*2 - 1.0
}

func (b *boid) calcAcceleration() {
	// occasionally change direction
	if rand.Float64() > 0.9999 {
		b.calcVelocity()
	}

	upper := b.position.AddV(viewRadius)
	lower := b.position.AddV(-viewRadius)
	avgPosition := util.Vector2D{}
	avgVelocity := util.Vector2D{}
	separation := util.Vector2D{}
	count := 0.0

	for i := math.Max(lower.X, 0); i <= math.Min(upper.X, screenWidth); i++ {
		for j := math.Max(lower.Y, 0); j <= math.Min(upper.Y, screenHeight); j++ {
			if id := boidMap[int(i)][int(j)]; id != -1 && id != b.id {
				if dist := boids[id].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[id].velocity)
					avgPosition = avgPosition.Add(boids[id].position)
					separation = separation.Add(b.position.Subtract(boids[id].position).DivisionV(dist))
				}
			}
		}
	}

	b.accel.X = b.borderBounce(b.position.X, screenWidth)
	b.accel.Y = b.borderBounce(b.position.Y, screenHeight)

	if count > 0 {
		avgPosition = avgPosition.DivisionV(count)
		avgVelocity = avgVelocity.DivisionV(count)
		accelAlignment := avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyV(adjRate)
		accelSeparation := separation.MultiplyV(adjRate)
		b.accel = b.accel.Add(accelAlignment).Add(accelCohesion).Add(accelSeparation)
	}
}

func (b *boid) borderBounce(pos, maxBorderPos float64) float64 {
	if pos < viewRadius {
		return 1 / pos
	} else if pos > maxBorderPos-viewRadius {
		return 1 / (pos - maxBorderPos)
	}

	return 0
}
