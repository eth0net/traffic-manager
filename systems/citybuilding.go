package systems

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Spritesheet contains city sprites.
var Spritesheet *common.Spritesheet

var cities = [...][12]int{
	{
		99, 100, 101,
		454, 269, 455,
		415, 195, 416,
		452, 306, 453,
	},
	{
		99, 100, 101,
		268, 269, 270,
		268, 269, 270,
		305, 306, 307,
	},
	{
		75, 76, 77,
		446, 261, 447,
		446, 261, 447,
		444, 298, 445,
	},
	{
		75, 76, 77,
		407, 187, 408,
		407, 187, 408,
		444, 298, 445,
	},
	{
		75, 76, 77,
		186, 150, 188,
		186, 150, 188,
		297, 191, 299,
	},
	{
		83, 84, 85,
		413, 228, 414,
		411, 191, 412,
		448, 302, 449,
	},
	{
		83, 84, 85,
		227, 228, 229,
		190, 191, 192,
		301, 302, 303,
	},
	{
		91, 92, 93,
		241, 242, 243,
		278, 279, 280,
		945, 946, 947,
	},
	{
		91, 92, 93,
		241, 242, 243,
		278, 279, 280,
		945, 803, 947,
	},
	{
		91, 92, 93,
		238, 239, 240,
		238, 239, 240,
		312, 313, 314,
	},
}

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
	usedTiles    []int
	elapsed      float32
}

// New is called to initialise the system when it is added to the scene.
func (cb *CityBuildingSystem) New(world *ecs.World) {
	fmt.Println("CityBuildingSystem was added to the Scene")

	rand.Seed(time.Now().UnixNano())

	// store world reference for later
	cb.world = world

	// load spritesheet
	Spritesheet = common.NewSpritesheetWithBorderFromFile(
		"textures/citySheet.png", 16, 16, 1, 1,
	)

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
	cb.elapsed += dt
	if cb.elapsed >= 10 {
		cb.generateCity()
		cb.elapsed = 0
	}
}

// generateCity generates a random city in a random map tile
func (cb *CityBuildingSystem) generateCity() {
	x := rand.Intn(18)
	y := rand.Intn(18)
	t := x + y*18

	for cb.isTileUsed(t) {
		if len(cb.usedTiles) > 300 {
			break // avoid infinite loop
		}
		x = rand.Intn(18)
		y = rand.Intn(18)
		t = x + y*18
	}
	cb.usedTiles = append(cb.usedTiles, t)

	city := rand.Intn(len(cities))
	cityTiles := []*City{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			tile := &City{BasicEntity: ecs.NewBasic()}
			tile.RenderComponent.Drawable = Spritesheet.Cell(cities[city][i+3*j])
			tile.RenderComponent.SetZIndex(1)
			tile.SpaceComponent.Position = engo.Point{
				X: float32(((x+1)*64)+8) + float32(i*16),
				Y: float32(((y + 1) * 64)) + float32(j*16),
			}
			cityTiles = append(cityTiles, tile)
		}
	}

	for _, system := range cb.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range cityTiles {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
}

func (cb *CityBuildingSystem) isTileUsed(tile int) bool {
	for _, t := range cb.usedTiles {
		if tile == t {
			return true
		}
	}
	return false
}
