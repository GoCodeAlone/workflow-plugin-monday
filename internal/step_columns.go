package internal

import (
	"context"
	"encoding/json"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type getColumnValuesStep struct{ name, moduleName string }

func newGetColumnValuesStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &getColumnValuesStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *getColumnValuesStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	if itemID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id is required"}}, nil
	}
	query := `query ($ids: [ID!]) { items(ids: $ids) { column_values { id text value } } }`
	var result struct {
		Items []struct {
			ColumnValues []any `json:"column_values"`
		} `json:"items"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"ids": []string{itemID}}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if len(result.Items) == 0 {
		return &sdk.StepResult{Output: map[string]any{"column_values": []any{}}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"column_values": result.Items[0].ColumnValues}}, nil
}

type changeColumnValueStep struct{ name, moduleName string }

func newChangeColumnValueStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &changeColumnValueStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *changeColumnValueStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	boardID := resolveValue("board_id", current, config)
	columnID := resolveValue("column_id", current, config)
	value := resolveValue("value", current, config)
	if itemID == "" || boardID == "" || columnID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id, board_id, and column_id are required"}}, nil
	}
	// value may be JSON string
	valMap := resolveMap("value", current, config)
	if valMap != nil {
		b, err := json.Marshal(valMap)
		if err != nil {
			return &sdk.StepResult{Output: map[string]any{"error": "marshal value: " + err.Error()}}, nil
		}
		value = string(b)
	}
	query := `mutation ($itemId: ID!, $boardId: ID!, $columnId: String!, $value: JSON!) {
		change_column_value(item_id: $itemId, board_id: $boardId, column_id: $columnId, value: $value) { id }
	}`
	var result struct {
		ChangeColumnValue map[string]any `json:"change_column_value"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"itemId": itemID, "boardId": boardID, "columnId": columnID, "value": value}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.ChangeColumnValue}, nil
}

type createColumnStep struct{ name, moduleName string }

func newCreateColumnStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createColumnStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createColumnStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	title := resolveValue("title", current, config)
	columnType := resolveValue("column_type", current, config)
	if boardID == "" || title == "" || columnType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id, title, and column_type are required"}}, nil
	}
	query := `mutation ($boardId: ID!, $title: String!, $columnType: ColumnType!) {
		create_column(board_id: $boardId, title: $title, column_type: $columnType) { id title type }
	}`
	var result struct {
		CreateColumn map[string]any `json:"create_column"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardId": boardID, "title": title, "columnType": columnType}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateColumn}, nil
}
