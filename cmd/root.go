/*
Copyright © 2022 Robert Greener <me@r0bert.dev>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the copyright holder nor the names of its contributors
    may be used to endorse or promote products derived from this software
    without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"encoding/json"
	"os"

	b "github.com/0xr0bert/gobelief/beliefspread"
	"github.com/0xr0bert/gobelief/runner"
	"github.com/google/uuid"
	"github.com/klauspost/compress/zstd"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gobelief",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := zap.NewProduction()
		if err != nil {
			return
		}
		config := new(runner.Configuration)

		startTime, err := cmd.Flags().GetUint32("start")

		if err != nil {
			logger.Error(
				"Failed to load start time",
				zap.String("errorMessage", err.Error()),
			)
			return
		}

		config.StartTime = b.SimTime(startTime)

		endTime, err := cmd.Flags().GetUint32("end")

		if err != nil {
			logger.Error(
				"Failed to load end time",
				zap.String("errorMessage", err.Error()),
			)
			return
		}

		config.EndTime = b.SimTime(endTime)

		outputFilepath, err := cmd.Flags().GetString("output")

		if err != nil {
			logger.Error(
				"Failed to get output filepath",
				zap.String("errorMessage", err.Error()),
			)

			return
		}

		if outputFilepath == "" {
			logger.Error(
				"outputFilepath is unset",
			)

			return
		}

		outputFile, err := os.Create(outputFilepath)

		if err != nil {
			logger.Error(
				"Failed to create output file",
				zap.String("errorMessage", err.Error()),
			)

			return
		}

		config.OutputFile = outputFile

		behavioursFilepath, err := cmd.Flags().GetString("behaviours")

		if err != nil {
			logger.Error(
				"Failed to get behaviours filepath",
				zap.String("errorMessage", err.Error()),
			)

			return
		}

		if behavioursFilepath == "" {
			logger.Error("behavioursFilepath is unset")

			return
		}

		behaviours, err := readBehavioursJson(behavioursFilepath)

		if err != nil {
			logger.Error(
				"Failed to read behaviours file",
				zap.String("errorMessage", err.Error()),
			)

			return
		}

		config.Behaviours = behaviours

		beliefsFilepath, err := cmd.Flags().GetString("beliefs")

		if err != nil {
			logger.Error(
				"Failed to get beliefs filepath",
				zap.String("errorMessage", err.Error()),
			)

			return
		}

		if beliefsFilepath == "" {
			logger.Error("beliefsFilepath unset")

			return
		}

		beliefs, err := readBeliefsJson(beliefsFilepath, behaviours)

		if err != nil {
			logger.Error(
				"Failed to read beliefs file",
				zap.String("errorMessage", err.Error()),
			)

			return
		}

		config.Beliefs = beliefs

		agentsFilepath, err := cmd.Flags().GetString("agents")

		if err != nil {
			logger.Error(
				"Failed to get agents filepath",
				zap.String("errorMessage", err.Error()),
			)

			return
		}

		if agentsFilepath == "" {
			logger.Error("agentsFilepath unset")

			return
		}

		agents, err := readAgentsJson(agentsFilepath, behaviours, beliefs)

		if err != nil {
			logger.Error(
				"Failed to read agents file",
				zap.String("errorMessage", err.Error()),
			)

			return
		}

		config.Agents = agents

		prsFilepath, err := cmd.Flags().GetString("prs")

		if err != nil {
			logger.Error(
				"Failed to get prs filepath",
				zap.String("errorMessage", err.Error()),
			)

			return
		}

		if prsFilepath == "" {
			logger.Error("prsFilepath unset")

			return
		}

		prs, err := readPrsJson(prsFilepath, beliefs, behaviours)

		if err != nil {
			logger.Error(
				"Failed to read prs file",
				zap.String("errorMessage", err.Error()),
			)

			return
		}

		config.Prs = prs

		simRunner := runner.Runner{
			Configuration: config,
			Logger:        logger,
		}
		simRunner.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Uint32P("start", "s", 1, "The start time of the simulation")
	rootCmd.Flags().Uint32P("end", "e", 1, "The end time of the simulation")
	rootCmd.Flags().StringP("output", "o", "", "The output file (e.g., output.json.zst)")
	rootCmd.Flags().StringP("behaviours", "b", "", "The behaviours.json file")
	rootCmd.Flags().StringP("beliefs", "c", "", "The beliefs.json file")
	rootCmd.Flags().StringP("agents", "a", "", "The agents.json.zst file")
	rootCmd.Flags().StringP("prs", "p", "", "The prs.json file")
}

func readBehavioursJson(path string) ([]*b.Behaviour, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	// var behaviourSpecs
	var behaviourSpecs []runner.BehaviourSpec
	err = json.Unmarshal(data, &behaviourSpecs)

	if err != nil {
		return nil, err
	}

	behaviours := make([]*b.Behaviour, len(behaviourSpecs))

	for i, spec := range behaviourSpecs {
		behaviours[i] = spec.ToBehaviour()
	}

	return behaviours, nil
}

func readBeliefsJson(path string, behaviours []*b.Behaviour) ([]*b.Belief, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var beliefSpecs []runner.BeliefSpec
	err = json.Unmarshal(data, &beliefSpecs)

	if err != nil {
		return nil, err
	}

	beliefs := make([]*b.Belief, len(beliefSpecs))

	for i, spec := range beliefSpecs {
		beliefs[i] = spec.ToBelief(behaviours)
	}

	for _, spec := range beliefSpecs {
		spec.LinkBeliefRelationships(beliefs)
	}

	return beliefs, nil
}

func readAgentsJson(
	path string,
	behaviours []*b.Behaviour,
	beliefs []*b.Belief,
) ([]*b.Agent, error) {
	decoder, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1))

	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	uncompressedData, err := decoder.DecodeAll(data, nil)

	data = nil

	if err != nil {
		return nil, err
	}

	var agentSpecs []runner.AgentSpec
	err = json.Unmarshal(uncompressedData, &agentSpecs)

	uncompressedData = nil

	if err != nil {
		return nil, err
	}

	agents := make([]*b.Agent, len(agentSpecs))

	for i, spec := range agentSpecs {
		agents[i] = spec.ToAgent(behaviours, beliefs)
	}

	uuidAgents := make(map[uuid.UUID]*b.Agent, len(agents))

	for _, agent := range agents {
		uuidAgents[agent.Uuid] = agent
	}

	for _, spec := range agentSpecs {
		spec.LinkFriends(uuidAgents)
	}

	return agents, nil
}

func readPrsJson(
	path string,
	beliefs []*b.Belief,
	behaviours []*b.Behaviour,
) (runner.PerformanceRelationships, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var specs []runner.PerformanceRelationshipSpec
	err = json.Unmarshal(data, &specs)

	uuidBeliefs := make(map[uuid.UUID]*b.Belief, len(beliefs))
	for _, belief := range beliefs {
		uuidBeliefs[belief.Uuid] = belief
	}

	uuidBehaviours := make(map[uuid.UUID]*b.Behaviour, len(behaviours))

	for _, behaviour := range behaviours {
		uuidBehaviours[behaviour.Uuid] = behaviour
	}

	return runner.PRSSpecToPerformanceRelationships(specs, uuidBeliefs, uuidBehaviours), nil
}
