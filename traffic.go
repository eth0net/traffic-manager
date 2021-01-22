package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/raziel2244/traffic-manager/systems"
)

type myScene struct{}

// Type identifies the scene type.
func (*myScene) Type() string {
	return "myGame"
}

// Preload is called before loading assets,
// allowing them to be registered / queued.
func (*myScene) Preload() {
	engo.Files.Load("textures/city.png")
}

// Setup is called before the main loop starts,
// allowing entities and systems to be added.
func (*myScene) Setup(u engo.Updater) {
	engo.Input.RegisterButton("AddCity", engo.KeyF1)
	common.SetBackground(color.White)

	var (
		scrollSpeed float32 = 400
		zoomSpeed   float32 = -0.125
		edgeMargin  float32 = 20
	)

	keyboardScroller := common.NewKeyboardScroller(
		scrollSpeed,
		engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis,
	)

	edgeScroller := &common.EdgeScroller{
		ScrollSpeed: scrollSpeed,
		EdgeMargin:  edgeMargin,
	}

	mouseZoomer := &common.MouseZoomer{
		ZoomSpeed: zoomSpeed,
	}

	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})

	world.AddSystem(keyboardScroller)
	world.AddSystem(edgeScroller)
	world.AddSystem(mouseZoomer)

	world.AddSystem(&systems.CityBuildingSystem{})
}

func main() {
	opts := engo.RunOptions{
		Title:          "Hello World",
		Width:          800,
		Height:         600,
		StandardInputs: true,
	}
	engo.Run(opts, &myScene{})
}
