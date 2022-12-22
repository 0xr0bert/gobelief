package runner

import (
	b "github.com/0xr0bert/gobelief/beliefspread"
	"github.com/google/uuid"
)

// PerformanceRelationships defines the relationship which defines how holding a belief affects the probability of
// performing a behaviour.
//
// The value should be in the range [-1,+1]
type PerformanceRelationships map[*b.Belief]map[*b.Behaviour]float64

// PRSSpecToPerformanceRelationships converts a slice of PerformanceRelationshipSpecs (i.e., what was read
// from JSON) to PerformanceRelationships.
//
// This takes the Beliefs (using a map from their UUID to the object) and
// Behaviours (using a map from their UUID to the object).
//
//goland:noinspection SpellCheckingInspection
func PRSSpecToPerformanceRelationships(
	prss []PerformanceRelationshipSpec,
	beliefs map[uuid.UUID]*b.Belief,
	behaviours map[uuid.UUID]*b.Behaviour,
) PerformanceRelationships {
	prs := make(PerformanceRelationships)
	for _, spec := range prss {
		belief := beliefs[spec.BeliefUuid]
		if belief != nil {
			_, found := prs[belief]
			if !found {
				prs[belief] = make(map[*b.Behaviour]float64)
			}

			behaviour := behaviours[spec.BehaviourUuid]
			if behaviour != nil {
				prs[belief][behaviour] = spec.Value
			}
		}
	}
	return prs
}
