package beliefspread_test

import (
	"testing"

	b "github.com/0xr0bert/gobelief/beliefspread"
)

func TestNewBehaviourAssignsRandomUuid(t *testing.T) {
	b1 := b.NewBehaviour("b1")
	b2 := b.NewBehaviour("b2")
	if b1.Uuid == b2.Uuid {
		t.Error("Equal UUIDs!")
	}
}

func TestNewBehaviourAssignsName(t *testing.T) {
	b := b.NewBehaviour("behaviour1")
	if b.Name != "behaviour1" {
		t.Errorf("Name should be behaviour1; it was %s", b.Name)
	}
}
