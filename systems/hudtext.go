package systems

import (
	"fmt"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// HUDTextMessageType is the entity type for HUDTextMessage.
const HUDTextMessageType string = "HUDTextMessage"

type mouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

// Text is an entity containing the printed text.
type Text struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// HUDTextMessage updates the HUD text
// with messages sent from other systems.
type HUDTextMessage struct {
	ecs.BasicEntity
	common.MouseComponent
	common.SpaceComponent
	Line1, Line2, Line3, Line4 string
}

// Type implements engo.Message Interface.
func (HUDTextMessage) Type() string {
	return HUDTextMessageType
}

// HUDMoneyMessageType is the type for HUDMoneyMessage.
const HUDMoneyMessageType string = "HUDMoneyMessage"

// HUDMoneyMessage updates the HUD text
// when changes occur with the player money.
type HUDMoneyMessage struct {
	Balance int
}

// Type implements the engo.Message interface.
func (HUDMoneyMessage) Type() string {
	return HUDMoneyMessageType
}

// HUDTextEntity is an entity for the HUDTextSystem.
// Keeps track of text position, size and contents.
type HUDTextEntity struct {
	*ecs.BasicEntity
	*common.MouseComponent
	*common.SpaceComponent
	Line1, Line2, Line3, Line4 string
}

// HUDTextSystem prints text to the
// screen based on current game state.
type HUDTextSystem struct {
	text1, text2, text3, text4, money Text

	entities []HUDTextEntity

	balance int
	updated bool
}

// New is called to initialise the system when it is
// added to the world. It loads the required resources,
// sets up entities and adss them to the world.
func (h *HUDTextSystem) New(w *ecs.World) {
	fnt := &common.Font{
		URL:  "go.ttf",
		FG:   color.Black,
		Size: 20,
	}
	fnt.CreatePreloaded()

	h.text1 = Text{BasicEntity: ecs.NewBasic()}
	h.text1.Drawable = common.Text{
		Font: fnt,
		Text: "Nothing selected!",
	}
	h.text1.SetShader(common.TextHUDShader)
	h.text1.SetZIndex(1001)
	h.text1.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: engo.WindowHeight() - 200},
	}

	h.text2 = Text{BasicEntity: ecs.NewBasic()}
	h.text2.Drawable = common.Text{
		Font: fnt,
		Text: "Click on an entity",
	}
	h.text2.SetShader(common.TextHUDShader)
	h.text2.SetZIndex(1001)
	h.text2.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: engo.WindowHeight() - 180},
	}

	h.text3 = Text{BasicEntity: ecs.NewBasic()}
	h.text3.Drawable = common.Text{
		Font: fnt,
		Text: "to get information",
	}
	h.text3.SetShader(common.TextHUDShader)
	h.text3.SetZIndex(1001)
	h.text3.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: engo.WindowHeight() - 160},
	}

	h.text4 = Text{BasicEntity: ecs.NewBasic()}
	h.text4.Drawable = common.Text{
		Font: fnt,
		Text: "about it here.",
	}
	h.text4.SetShader(common.TextHUDShader)
	h.text4.SetZIndex(1001)
	h.text4.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: engo.WindowHeight() - 140},
	}

	h.money = Text{BasicEntity: ecs.NewBasic()}
	h.money.Drawable = common.Text{
		Font: fnt,
		Text: "$0",
	}
	h.money.SetShader(common.TextHUDShader)
	h.money.SetZIndex(1001)
	h.money.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: engo.WindowHeight() - 40},
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(
				&h.text1.BasicEntity,
				&h.text1.RenderComponent,
				&h.text1.SpaceComponent,
			)
			sys.Add(
				&h.text2.BasicEntity,
				&h.text2.RenderComponent,
				&h.text2.SpaceComponent,
			)
			sys.Add(
				&h.text3.BasicEntity,
				&h.text3.RenderComponent,
				&h.text3.SpaceComponent,
			)
			sys.Add(
				&h.text4.BasicEntity,
				&h.text4.RenderComponent,
				&h.text4.SpaceComponent,
			)
			sys.Add(
				&h.money.BasicEntity,
				&h.money.RenderComponent,
				&h.money.SpaceComponent,
			)
		}
	}

	engo.Mailbox.Listen(HUDTextMessageType, func(msg engo.Message) {
		m, ok := msg.(HUDTextMessage)
		if !ok {
			return
		}
		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *common.MouseSystem:
				sys.Add(&m.BasicEntity, &m.MouseComponent, &m.SpaceComponent, nil)
			case *HUDTextSystem:
				sys.Add(
					&m.BasicEntity, &m.MouseComponent, &m.SpaceComponent,
					m.Line1, m.Line2, m.Line3, m.Line4,
				)
			}
		}
	})

	engo.Mailbox.Listen(HUDMoneyMessageType, func(msg engo.Message) {
		m, ok := msg.(HUDMoneyMessage)
		if !ok {
			return
		}
		h.balance = m.Balance
		h.updated = true
	})

	engo.Mailbox.Listen("WindowResizeMessage", func(msg engo.Message) {
		m, ok := msg.(engo.WindowResizeMessage)
		if !ok {
			return
		}

		heightNew := float32(m.NewHeight)
		h.text1.Position.Y = heightNew - 200
		h.text2.Position.Y = heightNew - 180
		h.text3.Position.Y = heightNew - 160
		h.text4.Position.Y = heightNew - 140
		h.money.Position.Y = heightNew - 40
	})
}

// Update is called for every frame, with dt set
// to the time in seconds since the last frame.
func (h *HUDTextSystem) Update(dt float32) {
	var text common.Text
	for _, entity := range h.entities {
		if entity.Clicked {
			text = h.text1.Drawable.(common.Text)
			text.Text = entity.Line1
			h.text1.Drawable = text
			text = h.text2.Drawable.(common.Text)
			text.Text = entity.Line2
			h.text2.Drawable = text
			text = h.text3.Drawable.(common.Text)
			text.Text = entity.Line3
			h.text3.Drawable = text
			text = h.text4.Drawable.(common.Text)
			text.Text = entity.Line4
			h.text4.Drawable = text
		}
	}

	if h.updated {
		text = h.money.Drawable.(common.Text)
		text.Text = fmt.Sprintf("$%v", h.balance)
		h.money.Drawable = text
	}
}

// Add takes an entity and adds it to the system.
func (h *HUDTextSystem) Add(
	b *ecs.BasicEntity,
	m *common.MouseComponent,
	s *common.SpaceComponent,
	l1, l2, l3, l4 string,
) {
	h.entities = append(h.entities, HUDTextEntity{b, m, s, l1, l2, l3, l4})
}

// Remove is called when an entity is removed from
// the world to remove it from the system as well.
func (h *HUDTextSystem) Remove(basic ecs.BasicEntity) {
	var del int = -1
	for i, entity := range h.entities {
		if entity.ID() == basic.ID() {
			del = i
			break
		}
	}
	if del >= 0 {
		h.entities = append(h.entities[:del], h.entities[del+1:]...)
	}
}
