package internal

import (
	"context"
	"encoding/json"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type createUpdateStep struct{ name, moduleName string }

func newCreateUpdateStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	body := resolveValue("body", current, config)
	if itemID == "" || body == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id and body are required"}}, nil
	}
	query := `mutation ($itemId: ID!, $body: String!) {
		create_update(item_id: $itemId, body: $body) { id body created_at }
	}`
	var result struct {
		CreateUpdate map[string]any `json:"create_update"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"itemId": itemID, "body": body}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateUpdate}, nil
}

type listUpdatesStep struct{ name, moduleName string }

func newListUpdatesStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listUpdatesStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listUpdatesStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	limit := resolveInt("limit", current, config)
	if limit == 0 {
		limit = 25
	}
	vars := map[string]any{"limit": limit}
	queryStr := `query ($limit: Int!) { updates(limit: $limit) { id body created_at } }`
	if itemID != "" {
		queryStr = `query ($ids: [ID!], $limit: Int!) { items(ids: $ids) { updates(limit: $limit) { id body created_at } } }`
		vars["ids"] = []string{itemID}
	}
	rawData, err := client.Execute(ctx, queryStr, vars)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	var result struct {
		Updates []any `json:"updates"`
	}
	if err := json.Unmarshal(rawData, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": "unmarshal: " + err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"updates": result.Updates}}, nil
}

type editUpdateStep struct{ name, moduleName string }

func newEditUpdateStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &editUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *editUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	updateID := resolveValue("update_id", current, config)
	body := resolveValue("body", current, config)
	if updateID == "" || body == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "update_id and body are required"}}, nil
	}
	query := `mutation ($id: ID!, $body: String!) {
		edit_update(id: $id, body: $body) { id body }
	}`
	var result struct {
		EditUpdate map[string]any `json:"edit_update"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"id": updateID, "body": body}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.EditUpdate}, nil
}

type deleteUpdateStep struct{ name, moduleName string }

func newDeleteUpdateStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &deleteUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	updateID := resolveValue("update_id", current, config)
	if updateID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "update_id is required"}}, nil
	}
	query := `mutation ($id: ID!) { delete_update(id: $id) { id } }`
	var result struct {
		DeleteUpdate map[string]any `json:"delete_update"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"id": updateID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.DeleteUpdate}, nil
}
