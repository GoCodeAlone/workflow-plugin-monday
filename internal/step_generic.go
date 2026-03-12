package internal

import (
	"context"
	"encoding/json"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type queryStep struct{ name, moduleName string }

func newQueryStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &queryStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *queryStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	if query == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "query is required"}}, nil
	}
	vars := resolveMap("variables", current, config)
	data, err := client.Execute(ctx, query, vars)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	var output map[string]any
	if err := json.Unmarshal(data, &output); err != nil {
		return &sdk.StepResult{Output: map[string]any{"data": string(data)}}, nil
	}
	return &sdk.StepResult{Output: output}, nil
}

type mutateStep struct{ name, moduleName string }

func newMutateStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &mutateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *mutateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	mutation := resolveValue("query", current, config)
	if mutation == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "query is required"}}, nil
	}
	vars := resolveMap("variables", current, config)
	data, err := client.Execute(ctx, mutation, vars)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	var output map[string]any
	if err := json.Unmarshal(data, &output); err != nil {
		return &sdk.StepResult{Output: map[string]any{"data": string(data)}}, nil
	}
	return &sdk.StepResult{Output: output}, nil
}
