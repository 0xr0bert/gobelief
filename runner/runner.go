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

type Configuration struct {
	Behaviours []*b.Behaviour
	Beliefs    []*b.Belief
	Agents     []*b.Agent
	Prs        PerformanceRelationships
	StartTime  b.SimTime
	EndTime    b.SimTime
	OutputFile *os.File
	FullOutput bool
}

type Runner struct {
	Configuration *Configuration
	Logger        *zap.Logger
}

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

func (r *Runner) tickBetween(start, end b.SimTime) {
	for i := start; i <= end; i++ {
		r.tick(i)
	}
}

func (r *Runner) tick(time b.SimTime) {
	r.Logger.Info("Perceiving beliefs", zap.Uint32("Day", uint32(time)))
	r.perceiveBeliefs(time)
	r.Logger.Info("Performing actions", zap.Uint32("Day", uint32(time)))
	r.performActions(time)
}

func (r *Runner) perceiveBeliefs(time b.SimTime) {
	for _, a := range r.Configuration.Agents {
		a.UpdateActivationForAllBeliefs(time, r.Configuration.Beliefs)
	}
}

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

func (r *Runner) performActions(time b.SimTime) {
	for _, a := range r.Configuration.Agents {
		r.agentPerformAction(a, time)
	}
}
