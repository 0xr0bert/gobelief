package beliefspread

import (
	"github.com/google/uuid"
)

// A Behaviour in the simulation.
type Behaviour struct {
	// The name of the behaviour.
	Name string
	// The UUID of the behaviour.
	Uuid uuid.UUID
}

// Create a new behaviour.
// Name is the name of the behaviour.
//
// This behaviour will have a randomly generated UUID.
func NewBehaviour(name string) (b *Behaviour) {
	b = new(Behaviour)
	b.Name = name
	b.Uuid, _ = uuid.NewRandom()

	return
}
