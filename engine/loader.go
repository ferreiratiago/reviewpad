// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/reviewpad/reviewpad/v3/handler"
	"gopkg.in/yaml.v3"
)

type LoadEnv struct {
	Visited map[string]bool
	Stack   map[string]bool
}

func hash(data []byte) string {
	dataHash := sha256.Sum256(data)
	dHash := fmt.Sprintf("%x", dataHash)
	return dHash
}

func Load(data []byte) (*ReviewpadFile, error) {
	file, err := parse(data)
	if err != nil {
		return nil, err
	}

	dHash := hash(data)

	visited := make(map[string]bool)
	stack := make(map[string]bool)
	visited[dHash] = true
	stack[dHash] = true

	env := &LoadEnv{
		Visited: visited,
		Stack:   stack,
	}

	file, err = processImports(file, env)
	if err != nil {
		return nil, err
	}

	file, err = normalize(file, inlineRulesNormalizer())
	if err != nil {
		return nil, err
	}

	return transform(file), nil
}

func parse(data []byte) (*ReviewpadFile, error) {
	file := ReviewpadFile{}
	err := yaml.Unmarshal([]byte(data), &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func transform(file *ReviewpadFile) *ReviewpadFile {
	var transformedRules []PadRule
	for _, rule := range file.Rules {
		kind := rule.Kind
		if rule.Kind == "" {
			kind = "patch"
		}
		transformedRules = append(transformedRules, PadRule{
			Name:        rule.Name,
			Kind:        kind,
			Description: rule.Description,
			Spec:        transformAladinoExpression(rule.Spec),
		})
	}

	var transformedWorkflows []PadWorkflow
	for _, workflow := range file.Workflows {
		var transformedRules []PadWorkflowRule
		for _, rule := range workflow.Rules {
			var transformedExtraActions []string
			for _, extraAction := range rule.ExtraActions {
				transformedExtraActions = append(transformedExtraActions, transformAladinoExpression(extraAction))
			}

			transformedRules = append(transformedRules, PadWorkflowRule{
				Rule:         rule.Rule,
				ExtraActions: transformedExtraActions,
			})
		}

		var transformedActions []string
		for _, action := range workflow.Actions {
			transformedActions = append(transformedActions, transformAladinoExpression(action))
		}

		transformedOn := []handler.TargetEntityKind{handler.PullRequest}
		if len(workflow.On) > 0 {
			transformedOn = workflow.On
		}

		transformedWorkflows = append(transformedWorkflows, PadWorkflow{
			Name:        workflow.Name,
			On:          transformedOn,
			Description: workflow.Description,
			Rules:       transformedRules,
			Actions:     transformedActions,
			AlwaysRun:   workflow.AlwaysRun,
		})
	}

	var transformedPipelines []PadPipeline

	for _, pipeline := range file.Pipelines {
		var transformedStages []PadStage

		for _, stage := range pipeline.Stages {
			var transformedActions []string
			for _, action := range stage.Actions {
				transformedActions = append(transformedActions, transformAladinoExpression(action))
			}

			transformedStages = append(transformedStages, PadStage{
				Actions: transformedActions,
				Until:   stage.Until,
			})
		}

		transformedPipelines = append(transformedPipelines, PadPipeline{
			Name:        pipeline.Name,
			Description: pipeline.Description,
			Trigger:     pipeline.Trigger,
			Stages:      transformedStages,
		})
	}

	return &ReviewpadFile{
		Version:      file.Version,
		Edition:      file.Edition,
		Mode:         file.Mode,
		IgnoreErrors: file.IgnoreErrors,
		Imports:      file.Imports,
		Groups:       file.Groups,
		Rules:        transformedRules,
		Labels:       file.Labels,
		Workflows:    transformedWorkflows,
		Pipelines:    transformedPipelines,
	}
}

func loadImport(reviewpadImport PadImport) (*ReviewpadFile, string, error) {
	resp, err := http.Get(reviewpadImport.Url)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	file, err := parse(content)
	if err != nil {
		return nil, "", err
	}

	return file, hash(content), nil
}

// processImports inlines the imports files into the current reviewpad file
// Post-condition: ReviewpadFile without import statements
func processImports(file *ReviewpadFile, env *LoadEnv) (*ReviewpadFile, error) {
	for _, reviewpadImport := range file.Imports {
		iFile, idHash, err := loadImport(reviewpadImport)
		if err != nil {
			return nil, err
		}

		// check for cycles
		if _, ok := env.Stack[idHash]; ok {
			return nil, fmt.Errorf("loader: cyclic dependency")
		}

		// optimize visits
		if _, ok := env.Visited[idHash]; ok {
			continue
		}

		// DFS call inline imports
		// update the environment
		env.Stack[idHash] = true
		env.Visited[idHash] = true

		subTreeFile, err := processImports(iFile, env)
		if err != nil {
			return nil, err
		}

		// remove from the stack
		delete(env.Stack, idHash)

		// append labels, rules and workflows
		file.appendLabels(subTreeFile)
		file.appendGroups(subTreeFile)
		file.appendRules(subTreeFile)
		file.appendWorkflows(subTreeFile)
	}

	// reset all imports
	file.Imports = []PadImport{}

	return file, nil
}
