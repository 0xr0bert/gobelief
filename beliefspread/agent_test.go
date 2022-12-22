package beliefspread

import (
	"math"
	"testing"
)

func TestNewAgentAssignsRandomUUID(t *testing.T) {
	a1 := NewAgent()
	a2 := NewAgent()
	if a1.Uuid == a2.Uuid {
		t.Error("Equal UUIDs!")
	}
}

func TestNewAgentAssignsActivationsEmpty(t *testing.T) {
	a := NewAgent()
	if len(a.Activations) != 0 {
		t.Error("Activations should be empty!")
	}
}

func TestNewAgentAssignsFriendsEmpty(t *testing.T) {
	a := NewAgent()
	if len(a.Friends) != 0 {
		t.Error("Friends should be empty!")
	}
}

func TestNewAgentAssignsActionsEmpty(t *testing.T) {
	a := NewAgent()
	if len(a.Actions) != 0 {
		t.Error("Actions should be empty!")
	}
}

func TestNewAgentAssignsDeltasEmpty(t *testing.T) {
	a := NewAgent()
	if len(a.Deltas) != 0 {
		t.Error("Deltas should be empty!")
	}
}

func TestWeightedRelationshipWhenExists(t *testing.T) {
	a := NewAgent()
	b1 := NewBelief("b1")
	b2 := NewBelief("b2")
	b1.Relationship[b2] = 0.5
	a.Activations[0] = make(map[*Belief]float64)
	a.Activations[0][b1] = 0.5
	wr := a.WeightedRelationship(0, b1, b2)
	if *wr != 0.25 {
		t.Errorf("Weighted relationship should be 0.25; it was %f", *wr)
	}
}

func TestWeightedRelationshipWhenActivationNotExists(t *testing.T) {
	a := NewAgent()
	b1 := NewBelief("b1")
	b2 := NewBelief("b2")
	b1.Relationship[b2] = 0.5
	wr := a.WeightedRelationship(0, b1, b2)
	if wr != nil {
		t.Errorf("Weighted relationship should be nil; it was %f", *wr)
	}
}

func TestWeightedRelationshipWhenRelationshipNotExists(t *testing.T) {
	a := NewAgent()
	b1 := NewBelief("b1")
	b2 := NewBelief("b2")
	a.Activations[0] = make(map[*Belief]float64)
	a.Activations[0][b1] = 0.5
	wr := a.WeightedRelationship(0, b1, b2)
	if wr != nil {
		t.Errorf("Weighted relationship should be nil; it was %f", *wr)
	}
}

func TestWeightedRelationshipWhenNotExists(t *testing.T) {
	a := NewAgent()
	b1 := NewBelief("b1")
	b2 := NewBelief("b2")
	wr := a.WeightedRelationship(0, b1, b2)
	if wr != nil {
		t.Errorf("Weighted relationship should be nil; it was %f", *wr)
	}
}

func TestWeightedRelationshipWhenTimeExistsButBeliefDoesntForActivation(t *testing.T) {
	a := NewAgent()
	b1 := NewBelief("b1")
	b2 := NewBelief("b2")
	b1.Relationship[b2] = 0.5
	a.Activations[0] = make(map[*Belief]float64)
	wr := a.WeightedRelationship(0, b1, b2)
	if wr != nil {
		t.Errorf("Weighted relationship should be nil; it was %f", *wr)
	}
}

func TestContextualiseWhenBeliefsEmptyReturns0(t *testing.T) {
	a := NewAgent()
	b1 := NewBelief("b1")
	b2 := NewBelief("b2")
	a.Activations[0] = make(map[*Belief]float64)
	a.Activations[0][b1] = 0.5
	a.Activations[0][b2] = 0.5
	c := a.Contextualise(0, b1, []*Belief{})
	if c != 0 {
		t.Errorf("Contextualise should be 0; it was %f", c)
	}
}

func TestContextualiseWhenBeliefsNonEmptyAndAllWeightedRelationshipsNotNil(t *testing.T) {
	a := NewAgent()
	b1 := NewBelief("b1")
	b2 := NewBelief("b2")
	b3 := NewBelief("b3")
	b1.Relationship[b2] = 0.5
	b1.Relationship[b3] = 0.5
	a.Activations[0] = make(map[*Belief]float64)
	a.Activations[0][b1] = 0.5
	a.Activations[0][b2] = 0.5
	a.Activations[0][b3] = 0.5
	c := a.Contextualise(0, b1, []*Belief{b2, b3})
	if c != 0.25 {
		t.Errorf("Contextualise should be 0.25; it was %f", c)
	}
}

func TestContextualiseWhenBeliefsNonEmptyAndSomeWeightedRelationshipsNil(t *testing.T) {
	a := NewAgent()
	b1 := NewBelief("b1")
	b2 := NewBelief("b2")
	b3 := NewBelief("b3")
	b1.Relationship[b2] = 0.5
	a.Activations[0] = make(map[*Belief]float64)
	a.Activations[0][b1] = 0.5
	a.Activations[0][b2] = 0.5
	a.Activations[0][b3] = 0.5
	c := a.Contextualise(0, b1, []*Belief{b2, b3})
	if c != 0.125 {
		t.Errorf("Contextualise should be 0.125; it was %f", c)
	}
}

func TestContextualiseWhenBeliefsNonEmptyAndAllWeightedRelationshipsNil(t *testing.T) {
	a := NewAgent()
	b1 := NewBelief("b1")
	b2 := NewBelief("b2")
	b3 := NewBelief("b3")
	a.Activations[0] = make(map[*Belief]float64)
	a.Activations[0][b1] = 0.5
	a.Activations[0][b2] = 0.5
	a.Activations[0][b3] = 0.5
	c := a.Contextualise(0, b1, []*Belief{b2, b3})
	if c != 0 {
		t.Errorf("Contextualise should be 0; it was %f", c)
	}
}

func TestGetActionsOfFriendsWhenFriendsEmpty(t *testing.T) {
	a := NewAgent()
	actions := a.GetActionsOfFriends(0)
	if len(actions) != 0 {
		t.Error("Actions should be empty!")
	}
}

func TestGetActionsOfFriendsWhenFriendsNotEmpty(t *testing.T) {
	a1 := NewAgent()
	a2 := NewAgent()
	a3 := NewAgent()
	a4 := NewAgent()
	a1.Friends[a1] = 0.2
	a1.Friends[a2] = 0.3
	a1.Friends[a3] = 0.5
	a1.Friends[a4] = 0.1

	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")

	a1.Actions[2] = b1
	a2.Actions[2] = b1
	a3.Actions[2] = b1
	a4.Actions[2] = b2

	actions := a1.GetActionsOfFriends(2)

	if len(actions) != 2 {
		t.Error("Actions should be length 2!")
	}

	if actions[b1] != 1.0 {
		t.Errorf("Actions should be 1.0; it was %f", actions[b1])
	}

	if actions[b2] != 0.1 {
		t.Errorf("Actions should be 0.1; it was %f", actions[b2])
	}
}

func TestPressureWhenNoFriends(t *testing.T) {
	a := NewAgent()
	bel := NewBelief("b")
	p := a.Pressure(bel, map[*Behaviour]float64{})
	if p != 0 {
		t.Errorf("Pressure should be 0; it was %f", p)
	}
}

func TestPressureWhenFriendsDidNothing(t *testing.T) {
	a1 := NewAgent()
	a2 := NewAgent()
	a3 := NewAgent()
	a4 := NewAgent()
	a1.Friends[a1] = 0.2
	a1.Friends[a2] = 0.3
	a1.Friends[a3] = 0.5
	a1.Friends[a4] = 0.1

	b1 := NewBelief("b1")

	p := a1.Pressure(b1, map[*Behaviour]float64{})

	if p != 0 {
		t.Errorf("Pressure should be 0; it was %f", p)
	}
}

func TestPressureWhenFriendsDidSomethingButPerceptionNil(t *testing.T) {
	agent := NewAgent()
	f1 := NewAgent()
	f2 := NewAgent()
	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")

	f1.Actions[2] = b1
	f2.Actions[2] = b2

	belief := NewBelief("b")

	agent.Friends[agent] = 0.2
	agent.Friends[f1] = 0.5
	agent.Friends[f2] = 1.0

	p := agent.Pressure(belief, agent.GetActionsOfFriends(2))

	if p != 0.0 {
		t.Errorf("Pressure should be 0.0; it was %f", p)
	}
}

func TestPressureWhenFriendsDidSomething(t *testing.T) {
	agent := NewAgent()
	f1 := NewAgent()
	f2 := NewAgent()

	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")

	f1.Actions[2] = b1
	f2.Actions[2] = b2

	belief := NewBelief("b")
	belief.Perception[b1] = 0.2
	belief.Perception[b2] = 0.3

	agent.Friends[f1] = 0.5
	agent.Friends[f2] = 1.0

	p := agent.Pressure(belief, agent.GetActionsOfFriends(2))

	if p != 0.2 {
		t.Errorf("Pressure should be 0.2; it was %f", p)
	}
}

func TestActivationChangeWhenPressurePositive(t *testing.T) {
	agent := NewAgent()
	f1 := NewAgent()
	f2 := NewAgent()

	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")

	f1.Actions[2] = b1
	f2.Actions[2] = b2

	belief := NewBelief("b")
	belief.Perception[b1] = 0.2
	belief.Perception[b2] = 0.3

	agent.Friends[f1] = 0.5
	agent.Friends[f2] = 1.0
	// Pressure is 0.2

	belief2 := NewBelief("b2")
	beliefs := []*Belief{belief, belief2}

	agent.Activations[2] = make(map[*Belief]float64)
	agent.Activations[2][belief] = 1.0
	agent.Activations[2][belief2] = 1.0

	belief.Relationship[belief] = 0.5
	belief.Relationship[belief2] = -0.75

	// Contextualise is -0.125

	change := agent.ActivationChange(2, belief, beliefs, agent.GetActionsOfFriends(2))

	if math.Abs(change-0.0875) > 0.000001 {
		t.Errorf("Change should be 0.0875; it was %f", change)
	}
}

func TestActivationChangeWhenPressureNegative(t *testing.T) {
	agent := NewAgent()
	f1 := NewAgent()
	f2 := NewAgent()

	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")

	f1.Actions[2] = b1
	f2.Actions[2] = b2

	belief := NewBelief("b")
	belief.Perception[b1] = -0.2
	belief.Perception[b2] = -0.3

	agent.Friends[f1] = 0.5
	agent.Friends[f2] = 1.0
	// Pressure is -0.2

	belief2 := NewBelief("b2")
	beliefs := []*Belief{belief, belief2}

	agent.Activations[2] = make(map[*Belief]float64)
	agent.Activations[2][belief] = 1.0
	agent.Activations[2][belief2] = 1.0

	belief.Relationship[belief] = 0.5
	belief.Relationship[belief2] = -0.75

	// Contextualise is -0.125

	change := agent.ActivationChange(2, belief, beliefs, agent.GetActionsOfFriends(2))

	if math.Abs(change-(-0.1125)) > 0.000001 {
		t.Errorf("Change should be -0.1125; it was %f", change)
	}
}

func TestMinWhenFirstSmaller(t *testing.T) {
	n1 := 0.2
	n2 := 0.5
	min := Min(n1, n2)
	if min != n1 {
		t.Errorf("Min should be %f; it was %f", n1, min)
	}
}

func TestMinWhenSecondSmaller(t *testing.T) {
	n1 := 0.5
	n2 := 0.2
	min := Min(n1, n2)
	if min != n2 {
		t.Errorf("Min should be %f; it was %f", n2, min)
	}
}

func TestMinWhenEqual(t *testing.T) {
	n1 := 0.2
	n2 := 0.2
	min := Min(n1, n2)
	if min != n2 {
		t.Errorf("Min should be %f; it was %f", n2, min)
	}
}

func TestMaxWhenFirstBigger(t *testing.T) {
	n1 := 0.5
	n2 := 0.2
	max := Max(n1, n2)
	if max != n1 {
		t.Errorf("Max should be %f; it was %f", n1, max)
	}
}

func TestMaxWhenSecondBigger(t *testing.T) {
	n1 := 0.2
	n2 := 0.5
	max := Max(n1, n2)
	if max != n2 {
		t.Errorf("Max should be %f; it was %f", n2, max)
	}
}

func TestUpdateActivationWhenPreviousActivationNone(t *testing.T) {
	agent := NewAgent()
	belief := NewBelief("belief1")
	beliefs := make([]*Belief, 0)

	agent.Deltas[belief] = 1.1

	expectedErrorText := "no activation for time"

	err := agent.UpdateActivation(3, belief, beliefs, agent.GetActionsOfFriends(3))

	if err == nil {
		t.Error("Expected error")
	}

	if err.Error() != expectedErrorText {
		t.Errorf("Expected error text %s; got %s", expectedErrorText, err.Error())
	}
}

func TestUpdateActivationWhenPreviousActivationAtTimeFoundButBeliefNot(t *testing.T) {
	agent := NewAgent()
	agent.Activations[2] = make(map[*Belief]float64)
	belief := NewBelief("belief1")
	beliefs := make([]*Belief, 0)

	agent.Deltas[belief] = 1.1

	expectedErrorText := "no activation found for belief"

	err := agent.UpdateActivation(3, belief, beliefs, agent.GetActionsOfFriends(3))

	if err == nil {
		t.Error("Expected error")
	}

	if err.Error() != expectedErrorText {
		t.Errorf("Expected error text %s; got %s", expectedErrorText, err.Error())
	}
}

func TestUpdateActivationWhenDeltaNone(t *testing.T) {
	agent := NewAgent()
	belief := NewBelief("belief1")
	beliefs := make([]*Belief, 0)

	expectedErrorText := "delta not found"

	err := agent.UpdateActivation(3, belief, beliefs, agent.GetActionsOfFriends(3))

	if err == nil {
		t.Error("Expected error")
	}

	if err.Error() != expectedErrorText {
		t.Errorf("Expected error text %s; got %s", expectedErrorText, err.Error())
	}
}

func TestUpdateActivationWhenNewValueInRange(t *testing.T) {
	agent := NewAgent()
	f1 := NewAgent()
	f2 := NewAgent()

	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")

	f1.Actions[2] = b1
	f2.Actions[2] = b2

	belief := NewBelief("b")
	belief.Perception[b1] = 0.2
	belief.Perception[b2] = 0.3
	agent.Friends[f1] = 0.5
	agent.Friends[f2] = 1.0

	// Pressure is 0.2

	belief2 := NewBelief("b2")
	beliefs := []*Belief{belief, belief2}

	agent.Activations[2] = make(map[*Belief]float64)
	agent.Activations[2][belief] = 0.5
	agent.Activations[2][belief2] = 1.0
	belief.Relationship[belief] = 1.0
	belief.Relationship[belief2] = -0.75

	// Contextualise is -0.0625

	// Activation change is 0.10625
	agent.Deltas[belief] = 1.1

	actions := agent.GetActionsOfFriends(2)

	err := agent.UpdateActivation(3, belief, beliefs, actions)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if agent.Activations[3][belief] != 0.65625 {
		t.Errorf("Activation should be 0.65625; it was %f", agent.Activations[3][belief])
	}
}

func TestUpdateActivationWhenNewValueTooLow(t *testing.T) {
	agent := NewAgent()
	f1 := NewAgent()
	f2 := NewAgent()

	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")

	f1.Actions[2] = b1
	f2.Actions[2] = b2

	belief := NewBelief("b")
	belief.Perception[b1] = 0.2
	belief.Perception[b2] = 0.3
	agent.Friends[f1] = 0.5
	agent.Friends[f2] = 1.0

	// Pressure is 0.2

	belief2 := NewBelief("b2")
	beliefs := []*Belief{belief, belief2}

	agent.Activations[2] = make(map[*Belief]float64)
	agent.Activations[2][belief] = 0.5
	agent.Activations[2][belief2] = 1.0
	belief.Relationship[belief] = 1.0
	belief.Relationship[belief2] = -0.75

	// Contextualise is -0.0625

	// Activation change is 0.10625
	agent.Deltas[belief] = -1000000

	// This is a total cheat to force activation really low, officially delta
	// cannot be less than 0, but it doesn't really matter.

	actions := agent.GetActionsOfFriends(2)

	err := agent.UpdateActivation(3, belief, beliefs, actions)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if agent.Activations[3][belief] != -1.0 {
		t.Errorf("Activation should be -1.0; it was %f", agent.Activations[3][belief])
	}
}

func TestUpdateActivationWhenNewValueTooHigh(t *testing.T) {
	agent := NewAgent()
	f1 := NewAgent()
	f2 := NewAgent()

	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")

	f1.Actions[2] = b1
	f2.Actions[2] = b2

	belief := NewBelief("b")
	belief.Perception[b1] = 0.2
	belief.Perception[b2] = 0.3
	agent.Friends[f1] = 0.5
	agent.Friends[f2] = 1.0

	// Pressure is 0.2

	belief2 := NewBelief("b2")
	beliefs := []*Belief{belief, belief2}

	agent.Activations[2] = make(map[*Belief]float64)
	agent.Activations[2][belief] = 0.5
	agent.Activations[2][belief2] = 1.0
	belief.Relationship[belief] = 1.0
	belief.Relationship[belief2] = -0.75

	// Contextualise is -0.0625

	// Activation change is 0.10625
	agent.Deltas[belief] = 1000000

	actions := agent.GetActionsOfFriends(2)

	err := agent.UpdateActivation(3, belief, beliefs, actions)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if agent.Activations[3][belief] != 1.0 {
		t.Errorf("Activation should be 1.0; it was %f", agent.Activations[3][belief])
	}
}

func TestUpdateActivationForAllBeliefsWhenNewValueInRange(t *testing.T) {
	agent := NewAgent()
	f1 := NewAgent()
	f2 := NewAgent()

	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")

	f1.Actions[2] = b1
	f2.Actions[2] = b2

	belief := NewBelief("b")
	belief.Perception[b1] = 0.2
	belief.Perception[b2] = 0.3
	agent.Friends[f1] = 0.5
	agent.Friends[f2] = 1.0

	// Pressure is 0.2

	belief2 := NewBelief("b2")
	beliefs := []*Belief{belief, belief2}

	agent.Activations[2] = make(map[*Belief]float64)
	agent.Activations[2][belief] = 0.5
	agent.Activations[2][belief2] = 1.0
	belief.Relationship[belief] = 1.0
	belief.Relationship[belief2] = -0.75

	// Contextualise is -0.0625

	// Activation change is 0.10625
	agent.Deltas[belief] = 1.1
	agent.Deltas[belief2] = 0.0

	err := agent.UpdateActivationForAllBeliefs(3, beliefs)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if agent.Activations[3][belief] != 0.65625 {
		t.Errorf("Activation should be 0.65625; it was %f", agent.Activations[3][belief])
	}
}

func TestUpdateActivationForAllBeliefsWhenErr(t *testing.T) {
	agent := NewAgent()
	f1 := NewAgent()
	f2 := NewAgent()

	b1 := NewBehaviour("b1")
	b2 := NewBehaviour("b2")

	f1.Actions[2] = b1
	f2.Actions[2] = b2

	belief := NewBelief("b")
	belief.Perception[b1] = 0.2
	belief.Perception[b2] = 0.3
	agent.Friends[f1] = 0.5
	agent.Friends[f2] = 1.0

	// Pressure is 0.2

	belief2 := NewBelief("b2")
	beliefs := []*Belief{belief, belief2}

	agent.Activations[2] = make(map[*Belief]float64)
	agent.Activations[2][belief] = 0.5
	agent.Activations[2][belief2] = 1.0
	belief.Relationship[belief] = 1.0
	belief.Relationship[belief2] = -0.75

	// Contextualise is -0.0625

	// Activation change is 0.10625
	agent.Deltas[belief] = 1.1

	err := agent.UpdateActivationForAllBeliefs(3, beliefs)
	if err == nil {
		t.Errorf("Expected error")
	}

	if err.Error() != "delta not found" {
		t.Errorf("Expected delta not found error; got %s", err.Error())
	}
}
