package main

import (
	"image"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/raziel2244/traffic-manager/systems"
)

// A HUD entity to display information on the
// screen that is positioned with the camera.
type HUD struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

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
		hudWidth    float32    = 200
		hudHeight   float32    = 200
		hudX        float32    = 0
		hudY        float32    = engo.WindowHeight() - hudHeight
		hudPosition engo.Point = engo.Point{X: hudX, Y: hudY}
		hudScale    engo.Point = engo.Point{X: 1, Y: 1}
	)

	hudImage := image.NewUniform(color.RGBA{205, 205, 205, 255})
	hudNRGBA := common.ImageToNRGBA(hudImage, int(hudWidth), int(hudHeight))
	hudImageObj := common.NewImageObject(hudNRGBA)
	hudTexture := common.NewTextureSingle(hudImageObj)

	hud := HUD{BasicEntity: ecs.NewBasic()}
	hud.RenderComponent = common.RenderComponent{
		Drawable: hudTexture,
		Scale:    hudScale,
		Repeat:   common.Repeat,
	}
	hud.SpaceComponent = common.SpaceComponent{
		Position: hudPosition,
		Width:    hudWidth,
		Height:   hudHeight,
	}
	hud.RenderComponent.SetShader(common.HUDShader)
	hud.RenderComponent.SetZIndex(1)

	engo.Mailbox.Listen("WindowResizeMessage", func(msg engo.Message) {
		resMsg, ok := msg.(engo.WindowResizeMessage)
		if !ok {
			return
		}

		hud.SpaceComponent.Position.Y = float32(resMsg.NewHeight) - hudHeight
	})

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

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hud.BasicEntity, &hud.RenderComponent, &hud.SpaceComponent)
		}
	}
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
