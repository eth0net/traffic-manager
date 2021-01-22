package systems

import (
	"fmt"
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

// A MouseTracker entity to keep track of the
// mouse position relative to the game grid.
type MouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

// CityBuildingSystem handles the
// creation of cities in the game.
type CityBuildingSystem struct {
	world        *ecs.World
	mouseTracker MouseTracker
}

// New is called to initialise the system when it is added to the scene.
func (cb *CityBuildingSystem) New(world *ecs.World) {
	fmt.Println("CityBuildingSystem was added to the Scene")

	// store world reference for later
	cb.world = world

	// prepare mouse tracker
	mt := &cb.mouseTracker
	mt.BasicEntity = ecs.NewBasic()
	mt.MouseComponent = common.MouseComponent{Track: true}

	// register the mouse tracker with relevant systems
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&mt.BasicEntity, &mt.MouseComponent, nil, nil)
		}
	}
}

// Remove is called when an entity is removed from
// the world to remove it from the system as well.
func (*CityBuildingSystem) Remove(ecs.BasicEntity) {}

// Update is called for every frame, with dt set
// to the time in seconds since the last frame.
func (cb *CityBuildingSystem) Update(dt float32) {
	// if the build city button has been pressed
	if engo.Input.Button("AddCity").JustPressed() {
		// load in the city texture
		texture, err := common.LoadedSprite("textures/city.png")
		if err != nil {
			log.Println("Unable to load texture: " + err.Error())
		}

		// set position to mouse position
		position := engo.Point{
			X: cb.mouseTracker.MouseX - 15,
			Y: cb.mouseTracker.MouseY - 48,
		}

		// create a city entity
		city := City{BasicEntity: ecs.NewBasic()}
		city.RenderComponent = common.RenderComponent{
			Scale:    engo.Point{X: .1, Y: .1},
			Drawable: texture,
		}
		city.SpaceComponent = common.SpaceComponent{
			Position: position,
			Width:    30,
			Height:   64,
		}

		// register the city with relevant systems
		for _, system := range cb.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&city.BasicEntity, &city.RenderComponent, &city.SpaceComponent)
			}
		}
	}
}
