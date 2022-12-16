package beliefspread

import (
	"errors"

	"github.com/google/uuid"
)

type Agent struct {
	Uuid        uuid.UUID
	Activations map[SimTime]map[*Belief]float64
	Friends     map[*Agent]float64
	Actions     map[SimTime]*Behaviour
	Deltas      map[*Belief]float64
}

func NewAgent() (a *Agent) {
	a = new(Agent)
	a.Uuid, _ = uuid.NewRandom()
	a.Activations = make(map[SimTime]map[*Belief]float64)
	a.Friends = make(map[*Agent]float64)
	a.Actions = make(map[SimTime]*Behaviour)
	a.Deltas = make(map[*Belief]float64)

	return
}

func (a *Agent) WeightedRelationship(t SimTime, b1 *Belief, b2 *Belief) *float64 {
	acts, found := a.Activations[t]

	if !found {
		return nil
	}

	b1Act, found := acts[b1]

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

func (a *Agent) Pressure(belief *Belief, actionsOfFriends map[*Behaviour]float64) (pressure float64) {
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

// Min returns the smaller of a or b.
func Min(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

// Max returns the larger of a or b.
func Max(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

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
