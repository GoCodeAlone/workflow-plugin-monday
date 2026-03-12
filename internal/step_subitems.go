package internal

import (
	"context"
	"encoding/json"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type createSubitemStep struct{ name, moduleName string }

func newCreateSubitemStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createSubitemStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createSubitemStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	parentItemID := resolveValue("parent_item_id", current, config)
	itemName := resolveValue("item_name", current, config)
	if parentItemID == "" || itemName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "parent_item_id and item_name are required"}}, nil
	}
	vars := map[string]any{"parentItemId": parentItemID, "itemName": itemName}
	colVals := resolveMap("column_values", current, config)
	if colVals != nil {
		b, err := json.Marshal(colVals)
		if err != nil {
			return &sdk.StepResult{Output: map[string]any{"error": "marshal column_values: " + err.Error()}}, nil
		}
		vars["columnValues"] = string(b)
	}
	query := `mutation ($parentItemId: ID!, $itemName: String!, $columnValues: JSON) {
		create_subitem(parent_item_id: $parentItemId, item_name: $itemName, column_values: $columnValues) { id name }
	}`
	var result struct {
		CreateSubitem map[string]any `json:"create_subitem"`
	}
	if err := client.ExecuteInto(ctx, query, vars, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateSubitem}, nil
}

type listSubitemsStep struct{ name, moduleName string }

func newListSubitemsStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listSubitemsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listSubitemsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	if itemID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id is required"}}, nil
	}
	query := `query ($ids: [ID!]) { items(ids: $ids) { subitems { id name state } } }`
	var result struct {
		Items []struct {
			Subitems []map[string]any `json:"subitems"`
		} `json:"items"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"ids": []string{itemID}}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if len(result.Items) == 0 {
		return &sdk.StepResult{Output: map[string]any{"subitems": []any{}}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"subitems": result.Items[0].Subitems}}, nil
}

type updateSubitemStep struct{ name, moduleName string }

func newUpdateSubitemStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &updateSubitemStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *updateSubitemStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	boardID := resolveValue("board_id", current, config)
	colVals := resolveMap("column_values", current, config)
	if itemID == "" || boardID == "" || colVals == nil {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id, board_id, and column_values are required"}}, nil
	}
	b, err := json.Marshal(colVals)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": "marshal column_values: " + err.Error()}}, nil
	}
	query := `mutation ($itemId: ID!, $boardId: ID!, $columnValues: JSON!) {
		change_multiple_column_values(item_id: $itemId, board_id: $boardId, column_values: $columnValues) { id name }
	}`
	var result struct {
		ChangeMultipleColumnValues map[string]any `json:"change_multiple_column_values"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"itemId": itemID, "boardId": boardID, "columnValues": string(b)}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.ChangeMultipleColumnValues}, nil
}

type deleteSubitemStep struct{ name, moduleName string }

func newDeleteSubitemStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &deleteSubitemStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteSubitemStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	if itemID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id is required"}}, nil
	}
	query := `mutation ($itemId: ID!) { delete_item(item_id: $itemId) { id } }`
	var result struct {
		DeleteItem map[string]any `json:"delete_item"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"itemId": itemID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.DeleteItem}, nil
}
