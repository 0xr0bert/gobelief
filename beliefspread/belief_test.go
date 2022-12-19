package beliefspread

import (
	"testing"
)

func TestNewBeliefAssignsRandomUuid(t *testing.T) {
	b1 := NewBelief("b1")
	b2 := NewBelief("b2")
	if b1.Uuid == b2.Uuid {
		t.Error("Equal UUIDs!")
	}
}

func TestNewBeliefAssignsName(t *testing.T) {
	b := NewBelief("belief1")
	if b.Name != "belief1" {
		t.Errorf("Name should be belief1; it was %s", b.Name)
	}
}
