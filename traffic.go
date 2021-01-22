package main

import (
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// A City entity within the game.
type City struct {
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
	common.SetBackground(color.White)

	cityTexture, err := common.LoadedSprite("textures/city.png")
	if err != nil {
		log.Println("Unable to load texture: " + err.Error())
	}

	city := City{
		BasicEntity: ecs.NewBasic(),
		RenderComponent: common.RenderComponent{
			Scale:    engo.Point{X: 1, Y: 1},
			Drawable: cityTexture,
		},
		SpaceComponent: common.SpaceComponent{
			Position: engo.Point{X: 10, Y: 10},
			Width:    303,
			Height:   641,
		},
	}

	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&city.BasicEntity, &city.RenderComponent, &city.SpaceComponent)
		}
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World",
		Width:  800,
		Height: 600,
	}
	engo.Run(opts, &myScene{})
}
