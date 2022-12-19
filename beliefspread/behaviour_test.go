package beliefspread

import (
	"testing"
)

func TestNewBehaviourAssignsRandomUuid(t *testing.T) {
	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")
	if b1.Uuid == b2.Uuid {
		t.Error("Equal UUIDs!")
	}
}

func TestNewBehaviourAssignsName(t *testing.T) {
	b := NewBehaviour("behaviour1")
	if b.Name != "behaviour1" {
		t.Errorf("Name should be behaviour1; it was %s", b.Name)
	}
}
