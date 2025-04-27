package camera

import (
	"image"
	"log"
	"math"
)

type State uint8

const (
	Inside State = iota
	Outside
)

type Camera struct {
	X, Y         float64
	state        State
	indoorBounds image.Rectangle
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

func (c *Camera) SetIndoorCameraBounds(bounds image.Rectangle) {
	log.Printf("set indoor camera bounds: x:%d, y:%d, max x: %d, max y: %d", bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)
	c.indoorBounds = bounds
}

func (c *Camera) UpdateState(state State) {
	c.state = state
}

func (c *Camera) State() State {
	return c.state
}

func (c *Camera) FollowTarget(targetX, targetY, screenWidth, screenHeight float64) {
	c.X = -targetX + screenWidth/2.0
	c.Y = -targetY + screenHeight/2.0
}

func (c *Camera) Constrain(TileMapWidthPixels, TileMapHeightPixels, screenWidth, screenHeight float64) {

	if c.State() == Outside {
		c.X = math.Min(c.X, 0)
		c.Y = math.Min(c.Y, -1)
		c.X = math.Max(c.X, screenWidth-TileMapWidthPixels)
		c.Y = math.Max(c.Y, screenHeight-TileMapHeightPixels)
	}

	if c.state == Inside {
		c.X = math.Min(c.X, -float64(c.indoorBounds.Min.X))
		c.Y = math.Min(c.Y, -float64(c.indoorBounds.Min.Y))
		c.X = math.Max(c.X, -float64(c.indoorBounds.Min.X))
		c.Y = math.Max(c.Y, screenHeight-float64(c.indoorBounds.Max.Y))

	}
}
