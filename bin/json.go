package main

import (
	"math"
	"sort"

	b "github.com/0xr0bert/gobelief/beliefspread"
	"github.com/google/uuid"
)

type BehaviourSpec struct {
	Name string    `json:"name"`
	Uuid uuid.UUID `json:"uuid"`
}

func (spec *BehaviourSpec) ToBehaviour() *b.Behaviour {
	behaviour := b.NewBehaviour(spec.Name)
	behaviour.Uuid = spec.Uuid
	return behaviour
}

type BeliefSpec struct {
	Name          string                `json:"name"`
	Uuid          uuid.UUID             `json:"uuid"`
	Perceptions   map[uuid.UUID]float64 `json:"perceptions"`
	Relationships map[uuid.UUID]float64 `json:"relationships"`
}

func (spec *BeliefSpec) ToBelief(behaviours []*b.Behaviour) *b.Belief {
	belief := b.NewBelief(spec.Name)
	belief.Uuid = spec.Uuid

	for _, behaviour := range behaviours {
		perception, found := spec.Perceptions[behaviour.Uuid]
		if found {
			belief.Perception[behaviour] = perception
		}
	}

	return belief
}

type PerformanceRelationshipSpec struct {
	BehaviourUuid uuid.UUID `json:"behaviourUuid"`
	BeliefUuid    uuid.UUID `json:"beliefUuid"`
	Value         float64   `json:"value"`
}

type AgentSpec struct {
	Uuid        uuid.UUID                           `json:"uuid"`
	Actions     map[b.SimTime]uuid.UUID             `json:"actions"`
	Activations map[b.SimTime]map[uuid.UUID]float64 `json:"activations"`
	Deltas      map[uuid.UUID]float64               `json:"deltas"`
	Friends     map[uuid.UUID]float64               `json:"friends"`
}

func (spec *AgentSpec) ToAgent(
	behaviours []*b.Behaviour,
	beliefs []*b.Belief,
) *b.Agent {
	a := b.NewAgent()
	a.Uuid = spec.Uuid

	uuidBehaviours := make(map[uuid.UUID]*b.Behaviour)
	for _, b := range behaviours {
		uuidBehaviours[b.Uuid] = b
	}

	for time, actionUuid := range spec.Actions {
		action := uuidBehaviours[actionUuid]
		if action != nil {
			a.Actions[time] = action
		}
	}

	uuidBeliefs := make(map[uuid.UUID]*b.Belief)
	for _, b := range beliefs {
		uuidBeliefs[b.Uuid] = b
	}

	for time, acts := range spec.Activations {
		a.Activations[time] = make(map[*b.Belief]float64)
		for beliefUuid, act := range acts {
			belief := uuidBeliefs[beliefUuid]
			if belief != nil {
				a.Activations[time][belief] = act
			}
		}
	}

	for beliefUuid, value := range spec.Deltas {
		belief := uuidBeliefs[beliefUuid]
		if belief != nil {
			a.Deltas[belief] = value
		}
	}

	return a
}

func (spec *AgentSpec) LinkFriends(agents map[uuid.UUID]*b.Agent) {
	thisAgent := agents[spec.Uuid]
	if thisAgent != nil {
		for friendUuid, w := range spec.Friends {
			friend := agents[friendUuid]
			if friend != nil {
				thisAgent.Friends[friend] = w
			}
		}
	}
}

type OutputSpec struct {
	MeanActivation         map[uuid.UUID]float64 `json:"meanActivation"`
	SDActivation           map[uuid.UUID]float64 `json:"sdActivation"`
	MedianActivation       map[uuid.UUID]float64 `json:"medianActivation"`
	NonzeroActivationCount map[uuid.UUID]uint64  `json:"nonzeroActivationCount"`
	NPerformers            map[uuid.UUID]uint64  `json:"nPerformers"`
}

func NewOutputSpec() *OutputSpec {
	o := new(OutputSpec)
	o.MeanActivation = make(map[uuid.UUID]float64)
	o.SDActivation = make(map[uuid.UUID]float64)
	o.MedianActivation = make(map[uuid.UUID]float64)
	o.NonzeroActivationCount = make(map[uuid.UUID]uint64)
	o.NPerformers = make(map[uuid.UUID]uint64)

	return o
}

type OutputSpecs struct {
	Data map[b.SimTime]OutputSpec `json:"data"`
}

func NewOutputSpecs(
	agents []*b.Agent,
	beliefs []*b.Belief,
	startTime b.SimTime,
	endTime b.SimTime,
) *OutputSpecs {
	data := make(map[b.SimTime]OutputSpec, endTime-startTime+1)

	for time := startTime; time <= endTime; time++ {
		o := NewOutputSpec()

		// Calculate mean activation
		for _, agent := range agents {
			acts := agent.Activations[time]
			for belief, act := range acts {
				o.MeanActivation[belief.Uuid] += act
			}
		}

		nAgents := len(agents)
		for u := range o.MeanActivation {
			o.MeanActivation[u] /= float64(nAgents)
		}

		// Calculate sd activation
		for _, agent := range agents {
			acts := agent.Activations[time]
			for belief, act := range acts {
				o.SDActivation[belief.Uuid] += math.Pow(
					act-o.MeanActivation[belief.Uuid],
					2.0,
				)
			}
		}

		for u, sd := range o.SDActivation {
			o.SDActivation[u] = math.Sqrt(sd / float64(nAgents-1))
		}

		// Calculate median activation
		activationsByUuid := make(map[uuid.UUID][]float64)

		for i, agent := range agents {
			for _, belief := range beliefs {
				_, found := activationsByUuid[belief.Uuid]
				if !found {
					activationsByUuid[belief.Uuid] = make([]float64, nAgents)
				}
				activationsByUuid[belief.Uuid][i] = agent.Activations[time][belief]
			}
		}

		middleIndex := nAgents / 2

		for uuid, acts := range activationsByUuid {
			sort.Float64s(acts)
			o.MedianActivation[uuid] = acts[middleIndex]
		}

		// Calculate non zero activation count
		for _, agent := range agents {
			for belief, activation := range agent.Activations[time] {
				if activation != 0.0 {
					o.NonzeroActivationCount[belief.Uuid]++
				}
			}
		}

		// Calculate n performers
		for _, agent := range agents {
			action := agent.Actions[time]
			if action != nil {
				o.NPerformers[action.Uuid]++
			}
		}

		data[time] = *o
	}

	return &OutputSpecs{Data: data}
}
