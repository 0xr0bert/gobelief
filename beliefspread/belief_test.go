package beliefspread_test

import (
	"testing"

	b "github.com/0xr0bert/gobelief/beliefspread"
)

func TestNewBeliefAssignsRandomUuid(t *testing.T) {
	b1 := b.NewBelief("b1")
	b2 := b.NewBelief("b2")
	if b1.Uuid == b2.Uuid {
		t.Error("Equal UUIDs!")
	}
}

func TestNewBeliefAssignsName(t *testing.T) {
	b := b.NewBelief("belief1")
	if b.Name != "belief1" {
		t.Errorf("Name should be belief1; it was %s", b.Name)
	}
}
