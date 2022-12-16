package runner

import (
	b "github.com/0xr0bert/gobelief/beliefspread"
	"github.com/google/uuid"
)

type PerformanceRelationships map[*b.Belief]map[*b.Behaviour]float64

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
