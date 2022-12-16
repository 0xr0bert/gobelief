package beliefspread

import (
	"github.com/google/uuid"
)

type Belief struct {
	Name         string
	Uuid         uuid.UUID
	Perception   map[*Behaviour]float64
	Relationship map[*Belief]float64
}

func NewBelief(name string) (b *Belief) {
	b = new(Belief)
	b.Name = name
	b.Uuid, _ = uuid.NewRandom()

	return
}
