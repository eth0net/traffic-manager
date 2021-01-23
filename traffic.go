package main

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/raziel2244/traffic-manager/systems"
	"golang.org/x/image/font/gofont/gosmallcaps"
)

const (
	cameraScrollSpeed float32 = 400
	cameraEdgeMargin  float32 = 20
	cameraZoomSpeed   float32 = -0.125
)

// A HUD entity to display information on the
// screen that is positioned with the camera.
type HUD struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// A Tile entity stores the contents
// of a single tile of the game map.
type Tile struct {
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
	engo.Files.Load("textures/citySheet.png", "tilemap/TrafficMap.tmx")
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
}

// Setup is called before the main loop starts,
// allowing entities and systems to be added.
func (*myScene) Setup(u engo.Updater) {
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
	hudTexture := common.NewTextureSingle(common.NewImageObject(hudNRGBA))

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
	hud.SetShader(common.HUDShader)
	hud.SetZIndex(1000)

	engo.Mailbox.Listen("WindowResizeMessage", func(msg engo.Message) {
		resMsg, ok := msg.(engo.WindowResizeMessage)
		if !ok {
			return
		}

		hud.SpaceComponent.Position.Y = float32(resMsg.NewHeight) - hudHeight
	})

	resource, err := engo.Files.Resource("tilemap/TrafficMap.tmx")
	if err != nil {
		panic(err)
	}
	levelData := resource.(common.TMXResource).Level
	common.CameraBounds = levelData.Bounds()

	tiles := []*Tile{}
	for _, tileLayer := range levelData.TileLayers {
		for _, tileElement := range tileLayer.Tiles {
			if tileElement.Image == nil {
				log.Printf("Tile is lacking image at point: %v", tileElement.Point)
			}
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			tile.RenderComponent = common.RenderComponent{
				Drawable: tileElement.Image,
				Scale:    engo.Point{X: 1, Y: 1},
			}
			tile.Position = tileElement.Point
			tiles = append(tiles, tile)
		}
	}

	keyboardScroller := common.NewKeyboardScroller(
		cameraScrollSpeed,
		engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis,
	)

	edgeScroller := &common.EdgeScroller{
		ScrollSpeed: cameraScrollSpeed,
		EdgeMargin:  cameraEdgeMargin,
	}

	mouseZoomer := &common.MouseZoomer{
		ZoomSpeed: cameraZoomSpeed,
	}

	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})

	world.AddSystem(keyboardScroller)
	world.AddSystem(edgeScroller)
	world.AddSystem(mouseZoomer)

	world.AddSystem(&systems.CityBuildingSystem{})
	world.AddSystem(&systems.HUDTextSystem{})

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hud.BasicEntity, &hud.RenderComponent, &hud.SpaceComponent)
			for _, tile := range tiles {
				sys.Add(&tile.BasicEntity, &tile.RenderComponent, &tile.SpaceComponent)
			}
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
