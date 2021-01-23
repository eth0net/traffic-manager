package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
)

// CityType indicates the type of a city.
type CityType int

const (
	// CityTypeNew is a brand new city.
	CityTypeNew CityType = iota
	// CityTypeTown is a town, level 1
	CityTypeTown
	// CityTypeCity is a city, level 2
	CityTypeCity
	// CityTypeMetro is a metro area, level 3
	CityTypeMetro
)

const (
	incomeTown     = 100
	incomeCity     = 500
	incomeMetro    = 1000
	incomeOfficers = -20
)

// CityUpdateMessageType is the entity type for CityUpdateMessage.
const CityUpdateMessageType string = "CityUpdateMessage"

// CityUpdateMessage updates cities from old to new.
type CityUpdateMessage struct {
	Old, New CityType
}

// Type implements engo.Message Interface.
func (CityUpdateMessage) Type() string {
	return CityUpdateMessageType
}

// AddOfficerMessageType is the entity type for AddOfficerMessage.
const AddOfficerMessageType string = "AddOfficerMessage"

// AddOfficerMessage adds an officer to the system.
type AddOfficerMessage struct{}

// Type implements engo.Message Interface.
func (AddOfficerMessage) Type() string {
	return AddOfficerMessageType
}

// MoneySystem keeps track of the players money.
type MoneySystem struct {
	balance               int
	towns, cities, metros int
	officers              int
	elapsed               float32
}

// New is called to initialise the system when it is added to the world.
// It sets up the various message listeners for the system.
func (m *MoneySystem) New(w *ecs.World) {
	engo.Mailbox.Listen(CityUpdateMessageType, func(msg engo.Message) {
		update, ok := msg.(CityUpdateMessage)
		if !ok {
			return
		}

		switch update.New {
		case CityTypeNew, CityTypeTown:
			m.towns++
		case CityTypeCity:
			m.cities++
		case CityTypeMetro:
			m.metros++
		}

		switch update.Old {
		case CityTypeTown:
			m.towns--
		case CityTypeCity:
			m.cities--
		case CityTypeMetro:
			m.metros--
		}
	})

	engo.Mailbox.Listen(AddOfficerMessageType, func(msg engo.Message) {
		m.officers++
	})
}

// Update is called every frame, with dt set to the time in
// seconds since the last frame.
//
// Money is regularly added for the number and type of cities
// connected and subtracted for the police force employed.
func (m *MoneySystem) Update(dt float32) {
	m.elapsed += dt
	if m.elapsed > 10 {
		m.balance += m.towns*incomeTown + m.cities*incomeCity +
			m.metros*incomeMetro + m.officers*incomeOfficers
		engo.Mailbox.Dispatch(HUDMoneyMessage{Balance: m.balance})
		m.elapsed = 0
	}
}

// Remove is called when an entity is removed from the world to remove it from the system as well.
// Does nothing since there are no entities in the system.
func (*MoneySystem) Remove(b ecs.BasicEntity) {}
