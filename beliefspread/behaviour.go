package beliefspread

import (
	"github.com/google/uuid"
)

type Behaviour struct {
	Name string
	Uuid uuid.UUID
}

func NewBehaviour(name string) (b *Behaviour) {
	b = new(Behaviour)
	b.Name = name
	b.Uuid, _ = uuid.NewRandom()

	return
}
