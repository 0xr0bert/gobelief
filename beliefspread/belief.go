package beliefspread

import (
	"github.com/google/uuid"
)

// A Belief in the simulation.
type Belief struct {
	// The name of the belief.
	Name string
	// The UUID of the belief.
	Uuid uuid.UUID
	// The perception of the belief to a behaviour.
	//
	// The perceiption is the amount an agent performing the behaviour can be
	// assumed to be driven by the belief.
	//
	// This should be a value between -1 and +1.
	Perception map[*Behaviour]float64

	// The relationship of the belief to another belief.
	//
	// The relationship is the amount the belief can be deemed to be compatible
	// with holding this belief, given that you already hold the other belief.
	Relationship map[*Belief]float64
}

// Create a new belief.
//
// This belief will have the supplied name, and a randomly generated UUID.
func NewBelief(name string) (b *Belief) {
	b = new(Belief)
	b.Name = name
	b.Uuid, _ = uuid.NewRandom()
	b.Perception = make(map[*Behaviour]float64)
	b.Relationship = make(map[*Belief]float64)

	return
}
