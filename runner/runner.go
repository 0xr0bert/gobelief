package runner

import (
	"encoding/json"
	"math/rand"
	"os"
	"sort"

	b "github.com/0xr0bert/gobelief/beliefspread"
	"github.com/klauspost/compress/zstd"
	"go.uber.org/zap"
)

// The configuration of the simulation.
type Configuration struct {
	// The behaviours in the simulation.
	Behaviours []*b.Behaviour
	// The beliefs in the simulation.
	Beliefs []*b.Belief
	// The agents in the simulation.
	Agents []*b.Agent
	// The PerformanceRelationships in the simulation.
	Prs PerformanceRelationships
	// The start time of the simulation.
	StartTime b.SimTime
	// The end time of the simulation (inclusive).
	EndTime b.SimTime
	// The output file.
	OutputFile *os.File
	// Whether to serialize the full state of agents, or just summary stats.
	FullOutput bool
}

// The simulation runner.
type Runner struct {
	// The configuration.
	Configuration *Configuration
	// The logger.
	Logger *zap.Logger
}

// Run the simulation.
func (r *Runner) Run() {
	r.Logger.Info(
		"Running simulation",
		zap.Uint32("Start", uint32(r.Configuration.StartTime)),
		zap.Uint32("End", uint32(r.Configuration.EndTime)),
		zap.Uint32("n beliefs", uint32(len(r.Configuration.Beliefs))),
		zap.Uint32("n behaviours", uint32(len(r.Configuration.Behaviours))),
		zap.Uint32("n agents", uint32(len(r.Configuration.Agents))),
	)
	r.tickBetween(r.Configuration.StartTime, r.Configuration.EndTime)
	r.Logger.Info("Ending simulation")
	var err error
	if r.Configuration.FullOutput {
		err = r.serializeFullOutput()
	} else {
		err = r.serializeOutput()
	}

	if err != nil {
		r.Logger.Error(
			"Error serializing output",
			zap.Error(err),
		)
	}
}

// Serialize the full state of agents as the output.
//
// This is stored as a zstd-compressed JSON file.
func (r *Runner) serializeFullOutput() error {
	specs := make([]*AgentSpec, len(r.Configuration.Agents))

	r.Logger.Info(
		"Preparing AgentSpecs for output",
	)

	for i, a := range r.Configuration.Agents {
		specs[i] = NewAgentSpecFromAgent(a)
	}

	r.Logger.Info(
		"Writing output to file",
		zap.String("File", r.Configuration.OutputFile.Name()),
	)

	data, err := json.Marshal(&specs)

	if err != nil {
		return err
	}

	zstdEncoder, err := zstd.NewWriter(r.Configuration.OutputFile)

	if err != nil {
		return err
	}

	encoder := json.NewEncoder(zstdEncoder)

	err = encoder.Encode(data)
	if err != nil {
		zstdEncoder.Close()
		return err
	}
	zstdEncoder.Close()
	return nil
}

// Serialize summary statistics about the agents.
//
// This calculates:
// - the number of agents performing each behaviour;
// - the mean activation for each belief;
// - the standard deviation of the activation for each belief;
// - the median activation for each belief; and
// - the number of agents who have non-zero activation for each belief.
//
// This is stored as a zstd-compressed JSON file.
func (r *Runner) serializeOutput() error {
	specs := NewOutputSpecs(
		r.Configuration.Agents,
		r.Configuration.Beliefs,
		r.Configuration.StartTime,
		r.Configuration.EndTime,
	)

	r.Logger.Info(
		"Writing output to file",
		zap.String("File", r.Configuration.OutputFile.Name()),
	)

	data, err := json.Marshal(&specs)

	if err != nil {
		return err
	}

	encoder, err := zstd.NewWriter(r.Configuration.OutputFile)

	if err != nil {
		return err
	}

	_, err = encoder.Write(data)
	if err != nil {
		encoder.Close()
		return err
	}
	encoder.Flush()
	encoder.Close()
	return nil
}

// Tick between two times (inclusive).
func (r *Runner) tickBetween(start, end b.SimTime) {
	for i := start; i <= end; i++ {
		r.tick(i)
	}
}

// "Tick" the simulation (run it for one time step - time).
func (r *Runner) tick(time b.SimTime) {
	r.Logger.Info("Perceiving beliefs", zap.Uint32("Day", uint32(time)))
	r.perceiveBeliefs(time)
	r.Logger.Info("Performing actions", zap.Uint32("Day", uint32(time)))
	r.performActions(time)
}

// Perceive the beliefs the agent holds for every agent.
//
// This updates all the agent's beliefs for every agent at the specified time
// step.
func (r *Runner) perceiveBeliefs(time b.SimTime) {
	for _, a := range r.Configuration.Agents {
		a.UpdateActivationForAllBeliefs(time, r.Configuration.Beliefs)
	}
}

// Perform an action for a specified agent at a specified time.
//
// If the preference for behaviours is fully negative, the "least-bad" option is
// chosen.
//
// If only one is positive, this option is chosen.
//
// If more than one is positive, it is chosen probabilistically based upon the
// preference.s
func (r *Runner) agentPerformAction(agent *b.Agent, time b.SimTime) {
	type probPair struct {
		behaviour *b.Behaviour
		value     float64
	}
	unnormalizedProbs := make([]probPair, len(r.Configuration.Behaviours))

	for i, b := range r.Configuration.Behaviours {
		unnormalizedProbs[i].behaviour = b
		for _, belief := range r.Configuration.Beliefs {
			prs := r.Configuration.Prs[belief][b]
			activation := agent.Activations[time][belief]
			unnormalizedProbs[i].value += prs * activation
		}
	}

	sort.Slice(unnormalizedProbs, func(i, j int) bool {
		return unnormalizedProbs[i].value < unnormalizedProbs[j].value
	})

	lastElem := unnormalizedProbs[len(unnormalizedProbs)-1]

	if lastElem.value < 0.0 {
		agent.Actions[time] = lastElem.behaviour
	} else {
		var filteredProbs []probPair
		for _, p := range unnormalizedProbs {
			if p.value > 0.0 {
				filteredProbs = append(filteredProbs, p)
			}
		}

		if len(filteredProbs) == 1 {
			agent.Actions[time] = filteredProbs[0].behaviour
		} else {
			normalizingFactor := 0.0
			for _, p := range filteredProbs {
				normalizingFactor += p.value
			}
			normalizedProbs := make([]probPair, len(filteredProbs))
			for i, p := range filteredProbs {
				normalizedProbs[i].behaviour = p.behaviour
				normalizedProbs[i].value = p.value / normalizingFactor
			}

			chosenBehaviour := normalizedProbs[len(normalizedProbs)-1].behaviour

			rv := rand.Float64()

			for _, p := range normalizedProbs {
				rv -= p.value
				if rv <= 0.0 {
					chosenBehaviour = p.behaviour
					break
				}
			}

			agent.Actions[time] = chosenBehaviour
		}
	}
}

// Perform actions for all agents at the specified time.
func (r *Runner) performActions(time b.SimTime) {
	for _, a := range r.Configuration.Agents {
		r.agentPerformAction(a, time)
	}
}
