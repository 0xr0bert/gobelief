package beliefspread

import (
	"errors"

	"github.com/google/uuid"
)

// An Agent in the simulation.
type Agent struct {
	// The UUID of the agent.
	Uuid uuid.UUID
	// The activation of the agent's beliefs at a given time.
	//
	// This should always be between -1 and +1.
	Activations map[SimTime]map[*Belief]float64
	// The relationship of the agent to other agents.
	//
	// This should be in the range [0, 1]
	Friends map[*Agent]float64
	// The actions of the agent at a given time.
	Actions map[SimTime]*Behaviour
	// The deltas of the agent's beliefs.
	//
	// This should be in the range [-1, +1] and is applied multiplicatively to
	// the activation of the belief at the next time step.
	Deltas map[*Belief]float64
}

// NewAgent creates a new agent with a randomly generated UUID.
func NewAgent() (a *Agent) {
	a = new(Agent)
	a.Uuid, _ = uuid.NewRandom()
	a.Activations = make(map[SimTime]map[*Belief]float64)
	a.Friends = make(map[*Agent]float64)
	a.Actions = make(map[SimTime]*Behaviour)
	a.Deltas = make(map[*Belief]float64)

	return
}

// WeightedRelationship gets the weighted relationship between two beliefs.
//
// This is the compatibility for holding b2, given that the Agent already holds
// b1.
//
// This is equal to the activation of b1 multiplied by the relationship between
// b1 and b2.
//
// Returns nil if the agent has no activation for b1, or if b1 and b2 have no
// relationship.
func (a *Agent) WeightedRelationship(t SimTime, b1 *Belief, b2 *Belief) *float64 {
	b1Act, found := a.Activations[t][b1]

	if !found {
		return nil
	}

	r, found := b1.Relationship[b2]

	if !found {
		return nil
	}

	returnVal := b1Act * r

	return &returnVal
}

// Contextualise gets the context for holding the Belief b
//
// This is the compatibility for holding b, given that the Agent all the beliefs
// the agent holds.
//
// This is an average of the weighted relationships for every Belief in beliefs.
func (a *Agent) Contextualise(t SimTime, b *Belief, beliefs []*Belief) (context float64) {
	size := len(beliefs)

	if size == 0 {
		return 0.0
	}

	for _, b2 := range beliefs {
		wr := a.WeightedRelationship(t, b, b2)

		if wr != nil {
			context += *wr
		}
	}

	context /= float64(size)

	return
}

// GetActionsOfFriends gets the actions of the agent's friends at a given time.
//
// The key is the behaviour, the value is the total weight of friends who
// performed that behaviour.
func (a *Agent) GetActionsOfFriends(t SimTime) (actions map[*Behaviour]float64) {
	actions = make(map[*Behaviour]float64)
	for friend, w := range a.Friends {
		action := friend.Actions[t]
		if action != nil {
			actions[action] += w
		}
	}

	return
}

// Pressure gets the pressure the Agent feels to adopt a Belief given the actions of
// their friends.
//
// This does not take into account the context of the Belief.
func (a *Agent) Pressure(
	belief *Belief,
	actionsOfFriends map[*Behaviour]float64,
) (pressure float64) {
	size := len(a.Friends)

	if size == 0 {
		return
	}

	for behaviour, w := range actionsOfFriends {
		pressure += belief.Perception[behaviour] * w
	}

	pressure /= float64(size)

	return
}

// ActivationChange gets the change in activation for the Agent as a result of observed
// Behaviour.
//
// This does take into account the context of the Belief.
func (a *Agent) ActivationChange(
	time SimTime,
	belief *Belief,
	beliefs []*Belief,
	actionsOfFriends map[*Behaviour]float64,
) float64 {
	pressure := a.Pressure(belief, actionsOfFriends)
	if pressure > 0.0 {
		return (1.0 + a.Contextualise(time, belief, beliefs)) / 2.0 * pressure
	} else {
		return (1.0 - a.Contextualise(time, belief, beliefs)) / 2.0 * pressure
	}
}

// Min returns the smallest of a or b.
func Min(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

// Max returns the largest of a or b.
func Max(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

// UpdateActivation updates the activation for a given belief and time, given the actions of
// the agents friends.
func (a *Agent) UpdateActivation(
	time SimTime,
	belief *Belief,
	beliefs []*Belief,
	actionsOfFriends map[*Behaviour]float64,
) error {
	delta, found := a.Deltas[belief]
	if !found {
		return errors.New("delta not found")
	}

	activations, found := a.Activations[time-1]

	if !found {
		return errors.New("no activation for time")
	}

	activation, found := activations[belief]

	if !found {
		return errors.New("no activation found for belief")
	}

	activationChange := a.ActivationChange(time-1, belief, beliefs, actionsOfFriends)

	newActivation := Max(-1.0, Min(1.0, delta*activation+activationChange))

	_, found = a.Activations[time]

	if !found {
		a.Activations[time] = make(map[*Belief]float64)
	}

	a.Activations[time][belief] = newActivation

	return nil
}

// UpdateActivationForAllBeliefs updates the activation for all beliefs at a given time.
func (a *Agent) UpdateActivationForAllBeliefs(
	time SimTime,
	beliefs []*Belief,
) error {
	actionsOfFriends := a.GetActionsOfFriends(time - 1)
	for _, belief := range beliefs {
		err := a.UpdateActivation(time, belief, beliefs, actionsOfFriends)
		if err != nil {
			return err
		}
	}

	return nil
}
