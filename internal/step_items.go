package internal

import (
	"context"
	"encoding/json"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type createItemStep struct{ name, moduleName string }

func newCreateItemStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createItemStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createItemStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	itemName := resolveValue("item_name", current, config)
	if boardID == "" || itemName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id and item_name are required"}}, nil
	}
	vars := map[string]any{"boardId": boardID, "itemName": itemName}
	groupID := resolveValue("group_id", current, config)
	if groupID != "" {
		vars["groupId"] = groupID
	}
	colVals := resolveMap("column_values", current, config)
	queryStr := `mutation ($boardId: ID!, $itemName: String!, $groupId: String, $columnValues: JSON) {
		create_item(board_id: $boardId, item_name: $itemName, group_id: $groupId, column_values: $columnValues) {
			id name state
		}
	}`
	if colVals != nil {
		b, err := json.Marshal(colVals)
		if err != nil {
			return &sdk.StepResult{Output: map[string]any{"error": "marshal column_values: " + err.Error()}}, nil
		}
		vars["columnValues"] = string(b)
	}
	var result struct {
		CreateItem map[string]any `json:"create_item"`
	}
	if err := client.ExecuteInto(ctx, queryStr, vars, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateItem}, nil
}

type listItemsStep struct{ name, moduleName string }

func newListItemsStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listItemsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listItemsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	if boardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id is required"}}, nil
	}
	limit := resolveInt("limit", current, config)
	if limit == 0 {
		limit = 25
	}
	query := `query ($ids: [ID!], $limit: Int!) {
		boards(ids: $ids) {
			items_page(limit: $limit) { cursor items { id name state } }
		}
	}`
	var result struct {
		Boards []struct {
			ItemsPage struct {
				Cursor string           `json:"cursor"`
				Items  []any `json:"items"`
			} `json:"items_page"`
		} `json:"boards"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"ids": []string{boardID}, "limit": limit}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if len(result.Boards) == 0 {
		return &sdk.StepResult{Output: map[string]any{"items": []any{}}}, nil
	}
	page := result.Boards[0].ItemsPage
	return &sdk.StepResult{Output: map[string]any{"items": page.Items, "cursor": page.Cursor}}, nil
}

type fetchItemStep struct{ name, moduleName string }

func newFetchItemStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &fetchItemStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *fetchItemStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	if itemID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id is required"}}, nil
	}
	query := `query ($ids: [ID!]) { items(ids: $ids) { id name state board { id name } } }`
	var result struct {
		Items []any `json:"items"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"ids": []string{itemID}}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if len(result.Items) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "item not found"}}, nil
	}
	return &sdk.StepResult{Output: result.Items[0].(map[string]any)}, nil
}

type updateItemStep struct{ name, moduleName string }

func newUpdateItemStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &updateItemStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *updateItemStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	boardID := resolveValue("board_id", current, config)
	if itemID == "" || boardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id and board_id are required"}}, nil
	}
	colVals := resolveMap("column_values", current, config)
	if colVals == nil {
		return &sdk.StepResult{Output: map[string]any{"error": "column_values is required"}}, nil
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

type moveItemStep struct{ name, moduleName string }

func newMoveItemStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &moveItemStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *moveItemStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	groupID := resolveValue("group_id", current, config)
	if itemID == "" || groupID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id and group_id are required"}}, nil
	}
	query := `mutation ($itemId: ID!, $groupId: String!) {
		move_item_to_group(item_id: $itemId, group_id: $groupId) { id }
	}`
	var result struct {
		MoveItemToGroup map[string]any `json:"move_item_to_group"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"itemId": itemID, "groupId": groupID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.MoveItemToGroup}, nil
}

type archiveItemStep struct{ name, moduleName string }

func newArchiveItemStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &archiveItemStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *archiveItemStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	if itemID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id is required"}}, nil
	}
	query := `mutation ($itemId: ID!) { archive_item(item_id: $itemId) { id state } }`
	var result struct {
		ArchiveItem map[string]any `json:"archive_item"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"itemId": itemID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.ArchiveItem}, nil
}

type deleteItemStep struct{ name, moduleName string }

func newDeleteItemStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &deleteItemStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteItemStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
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

type searchItemsStep struct{ name, moduleName string }

func newSearchItemsStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &searchItemsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *searchItemsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	term := resolveValue("term", current, config)
	if term == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "term is required"}}, nil
	}
	limit := resolveInt("limit", current, config)
	if limit == 0 {
		limit = 25
	}
	query := `query ($term: String!, $limit: Int!) {
		items_by_multiple_column_values(limit: $limit, column_id: "name", column_value: $term) { id name state }
	}`
	var result struct {
		Items []any `json:"items_by_multiple_column_values"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"term": term, "limit": limit}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"items": result.Items}}, nil
}
