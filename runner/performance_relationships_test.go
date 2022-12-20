package runner

import (
	"testing"

	b "github.com/0xr0bert/gobelief/beliefspread"
	"github.com/google/uuid"
)

func TestPrsSpecToPerformanceRelationshipsWhenPrssEmpty(t *testing.T) {
	bel1 := b.NewBelief("bel1")
	bel2 := b.NewBelief("bel2")

	beh1 := b.NewBehaviour("beh1")
	beh2 := b.NewBehaviour("beh2")

	prss := make([]PerformanceRelationshipSpec, 4)

	beliefs := map[uuid.UUID]*b.Belief{
		bel1.Uuid: bel1,
		bel2.Uuid: bel2,
	}

	behaviours := map[uuid.UUID]*b.Behaviour{
		beh1.Uuid: beh1,
		beh2.Uuid: beh2,
	}

	prs := PRSSpecToPerformanceRelationships(prss, beliefs, behaviours)

	if len(prs) != 0 {
		t.Errorf("len(prs) should be 0; it was %d", len(prs))
	}
}

func TestPRSSpecToPerformanceRelationshipsWhenAllOK(t *testing.T) {
	bel1 := b.NewBelief("bel1")
	bel2 := b.NewBelief("bel2")

	beh1 := b.NewBehaviour("beh1")
	beh2 := b.NewBehaviour("beh2")

	prss := make([]PerformanceRelationshipSpec, 4)
	prss[0] = PerformanceRelationshipSpec{
		BeliefUuid:    bel1.Uuid,
		BehaviourUuid: beh1.Uuid,
		Value:         0.2,
	}
	prss[1] = PerformanceRelationshipSpec{
		BeliefUuid:    bel1.Uuid,
		BehaviourUuid: beh2.Uuid,
		Value:         -0.2,
	}
	prss[2] = PerformanceRelationshipSpec{
		BeliefUuid:    bel2.Uuid,
		BehaviourUuid: beh1.Uuid,
		Value:         0.5,
	}
	prss[3] = PerformanceRelationshipSpec{
		BeliefUuid:    bel2.Uuid,
		BehaviourUuid: beh2.Uuid,
		Value:         -0.6,
	}

	beliefs := map[uuid.UUID]*b.Belief{
		bel1.Uuid: bel1,
		bel2.Uuid: bel2,
	}

	behaviours := map[uuid.UUID]*b.Behaviour{
		beh1.Uuid: beh1,
		beh2.Uuid: beh2,
	}

	prs := PRSSpecToPerformanceRelationships(prss, beliefs, behaviours)

	if prs[bel1][beh1] != 0.2 {
		t.Errorf("prs[bel1][beh1] should have been 0.2; it was %f", prs[bel1][beh1])
	}

	if prs[bel1][beh2] != -0.2 {
		t.Errorf("prs[bel1][beh2] should have been -0.2; it was %f", prs[bel1][beh2])
	}

	if prs[bel2][beh1] != 0.5 {
		t.Errorf("prs[bel2][beh1] should have been 0.5; it was %f", prs[bel2][beh1])
	}

	if prs[bel2][beh2] != -0.6 {
		t.Errorf("prs[bel2][beh2] should have been -0.6; it was %f", prs[bel2][beh2])
	}

	if len(prs) != 2 {
		t.Errorf("len(prs) should have been 2; it was %d", len(prs))
	}

	if len(prs[bel1]) != 2 {
		t.Errorf("len(prs[bel1]) should have been 2; it was %d", len(prs[bel1]))
	}

	if len(prs[bel2]) != 2 {
		t.Errorf("len(prs[bel2]) should have been 2; it was %d", len(prs[bel2]))
	}
}

func TestPRSSpecToPerformanceRelationshipsWhenSomeBeliefsMissing(t *testing.T) {
	bel1 := b.NewBelief("bel1")
	bel2 := b.NewBelief("bel2")

	beh1 := b.NewBehaviour("beh1")
	beh2 := b.NewBehaviour("beh2")

	prss := make([]PerformanceRelationshipSpec, 4)
	prss[0] = PerformanceRelationshipSpec{
		BeliefUuid:    bel1.Uuid,
		BehaviourUuid: beh1.Uuid,
		Value:         0.2,
	}
	prss[1] = PerformanceRelationshipSpec{
		BeliefUuid:    bel1.Uuid,
		BehaviourUuid: beh2.Uuid,
		Value:         -0.2,
	}
	prss[2] = PerformanceRelationshipSpec{
		BeliefUuid:    bel2.Uuid,
		BehaviourUuid: beh1.Uuid,
		Value:         0.5,
	}
	prss[3] = PerformanceRelationshipSpec{
		BeliefUuid:    bel2.Uuid,
		BehaviourUuid: beh2.Uuid,
		Value:         -0.6,
	}

	beliefs := map[uuid.UUID]*b.Belief{
		bel1.Uuid: bel1,
	}

	behaviours := map[uuid.UUID]*b.Behaviour{
		beh1.Uuid: beh1,
		beh2.Uuid: beh2,
	}

	prs := PRSSpecToPerformanceRelationships(prss, beliefs, behaviours)

	if prs[bel1][beh1] != 0.2 {
		t.Errorf("prs[bel1][beh1] should have been 0.2; it was %f", prs[bel1][beh1])
	}

	if prs[bel1][beh2] != -0.2 {
		t.Errorf("prs[bel1][beh2] should have been -0.2; it was %f", prs[bel1][beh2])
	}

	if len(prs) != 1 {
		t.Errorf("len(prs) should have been 2; it was %d", len(prs))
	}

	if len(prs[bel1]) != 2 {
		t.Errorf("len(prs[bel1]) should have been 2; it was %d", len(prs[bel1]))
	}
}

func TestPRSSpecToPerformanceRelationshipsWhenSomeBehavioursMissing(t *testing.T) {
	bel1 := b.NewBelief("bel1")
	bel2 := b.NewBelief("bel2")

	beh1 := b.NewBehaviour("beh1")
	beh2 := b.NewBehaviour("beh2")

	prss := make([]PerformanceRelationshipSpec, 4)
	prss[0] = PerformanceRelationshipSpec{
		BeliefUuid:    bel1.Uuid,
		BehaviourUuid: beh1.Uuid,
		Value:         0.2,
	}
	prss[1] = PerformanceRelationshipSpec{
		BeliefUuid:    bel1.Uuid,
		BehaviourUuid: beh2.Uuid,
		Value:         -0.2,
	}
	prss[2] = PerformanceRelationshipSpec{
		BeliefUuid:    bel2.Uuid,
		BehaviourUuid: beh1.Uuid,
		Value:         0.5,
	}
	prss[3] = PerformanceRelationshipSpec{
		BeliefUuid:    bel2.Uuid,
		BehaviourUuid: beh2.Uuid,
		Value:         -0.6,
	}

	beliefs := map[uuid.UUID]*b.Belief{
		bel1.Uuid: bel1,
		bel2.Uuid: bel2,
	}

	behaviours := map[uuid.UUID]*b.Behaviour{
		beh1.Uuid: beh1,
	}

	prs := PRSSpecToPerformanceRelationships(prss, beliefs, behaviours)

	if prs[bel1][beh1] != 0.2 {
		t.Errorf("prs[bel1][beh1] should have been 0.2; it was %f", prs[bel1][beh1])
	}

	if prs[bel2][beh1] != 0.5 {
		t.Errorf("prs[bel2][beh1] should have been 0.5; it was %f", prs[bel2][beh1])
	}

	if len(prs) != 2 {
		t.Errorf("len(prs) should have been 2; it was %d", len(prs))
	}

	if len(prs[bel1]) != 1 {
		t.Errorf("len(prs[bel1]) should have been 1; it was %d", len(prs[bel1]))
	}

	if len(prs[bel2]) != 1 {
		t.Errorf("len(prs[bel2]) should have been 1; it was %d", len(prs[bel2]))
	}
}

func TestPRSSpecToPerformanceRelationshipsWhenSomeBehavioursAndBeliefsMissing(t *testing.T) {
	bel1 := b.NewBelief("bel1")
	bel2 := b.NewBelief("bel2")

	beh1 := b.NewBehaviour("beh1")
	beh2 := b.NewBehaviour("beh2")

	prss := make([]PerformanceRelationshipSpec, 4)
	prss[0] = PerformanceRelationshipSpec{
		BeliefUuid:    bel1.Uuid,
		BehaviourUuid: beh1.Uuid,
		Value:         0.2,
	}
	prss[1] = PerformanceRelationshipSpec{
		BeliefUuid:    bel1.Uuid,
		BehaviourUuid: beh2.Uuid,
		Value:         -0.2,
	}
	prss[2] = PerformanceRelationshipSpec{
		BeliefUuid:    bel2.Uuid,
		BehaviourUuid: beh1.Uuid,
		Value:         0.5,
	}
	prss[3] = PerformanceRelationshipSpec{
		BeliefUuid:    bel2.Uuid,
		BehaviourUuid: beh2.Uuid,
		Value:         -0.6,
	}

	beliefs := map[uuid.UUID]*b.Belief{
		bel1.Uuid: bel1,
	}

	behaviours := map[uuid.UUID]*b.Behaviour{
		beh1.Uuid: beh1,
	}

	prs := PRSSpecToPerformanceRelationships(prss, beliefs, behaviours)

	if prs[bel1][beh1] != 0.2 {
		t.Errorf("prs[bel1][beh1] should have been 0.2; it was %f", prs[bel1][beh1])
	}

	if len(prs) != 1 {
		t.Errorf("len(prs) should have been 1; it was %d", len(prs))
	}

	if len(prs[bel1]) != 1 {
		t.Errorf("len(prs[bel1]) should have been 1; it was %d", len(prs[bel1]))
	}
}
