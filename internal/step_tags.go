package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type listTagsStep struct{ name, moduleName string }

func newListTagsStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listTagsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listTagsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, _ map[string]any, _ map[string]any, _ map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	query := `query { tags { id name color } }`
	var result struct {
		Tags []any `json:"tags"`
	}
	if err := client.ExecuteInto(ctx, query, nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"tags": result.Tags}}, nil
}

type createTagStep struct{ name, moduleName string }

func newCreateTagStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createTagStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createTagStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	tagName := resolveValue("tag_name", current, config)
	if tagName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "tag_name is required"}}, nil
	}
	query := `mutation ($tagName: String!) {
		create_or_get_tag(tag_name: $tagName) { id name color }
	}`
	var result struct {
		CreateOrGetTag map[string]any `json:"create_or_get_tag"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"tagName": tagName}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateOrGetTag}, nil
}
