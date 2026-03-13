package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type createGroupStep struct{ name, moduleName string }

func newCreateGroupStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createGroupStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createGroupStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	groupName := resolveValue("group_name", current, config)
	if boardID == "" || groupName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id and group_name are required"}}, nil
	}
	query := `mutation ($boardId: ID!, $groupName: String!) {
		create_group(board_id: $boardId, group_name: $groupName) { id title }
	}`
	var result struct {
		CreateGroup map[string]any `json:"create_group"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardId": boardID, "groupName": groupName}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateGroup}, nil
}

type listGroupsStep struct{ name, moduleName string }

func newListGroupsStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listGroupsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listGroupsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	if boardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id is required"}}, nil
	}
	query := `query ($ids: [ID!]) { boards(ids: $ids) { groups { id title color } } }`
	var result struct {
		Boards []struct {
			Groups []any `json:"groups"`
		} `json:"boards"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"ids": []string{boardID}}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if len(result.Boards) == 0 {
		return &sdk.StepResult{Output: map[string]any{"groups": []any{}}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"groups": result.Boards[0].Groups}}, nil
}

type updateGroupStep struct{ name, moduleName string }

func newUpdateGroupStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &updateGroupStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *updateGroupStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	groupID := resolveValue("group_id", current, config)
	attribute := resolveValue("attribute", current, config)
	newValue := resolveValue("new_value", current, config)
	if boardID == "" || groupID == "" || attribute == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id, group_id, and attribute are required"}}, nil
	}
	query := `mutation ($boardId: ID!, $groupId: String!, $groupAttribute: GroupAttributes!, $newValue: String!) {
		update_group(board_id: $boardId, group_id: $groupId, group_attribute: $groupAttribute, new_value: $newValue) { id title }
	}`
	var result struct {
		UpdateGroup map[string]any `json:"update_group"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardId": boardID, "groupId": groupID, "groupAttribute": attribute, "newValue": newValue}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.UpdateGroup}, nil
}

type moveGroupStep struct{ name, moduleName string }

func newMoveGroupStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &moveGroupStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *moveGroupStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	groupID := resolveValue("group_id", current, config)
	targetBoardID := resolveValue("target_board_id", current, config)
	if boardID == "" || groupID == "" || targetBoardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id, group_id, and target_board_id are required"}}, nil
	}
	query := `mutation ($boardId: ID!, $groupId: String!, $targetBoardId: ID!) {
		move_items_to_board(board_id: $boardId, group_id: $groupId, target_board_id: $targetBoardId) { id }
	}`
	var result struct {
		MoveItemsToBoard []any `json:"move_items_to_board"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardId": boardID, "groupId": groupID, "targetBoardId": targetBoardID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"items": result.MoveItemsToBoard}}, nil
}

type deleteGroupStep struct{ name, moduleName string }

func newDeleteGroupStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &deleteGroupStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteGroupStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	groupID := resolveValue("group_id", current, config)
	if boardID == "" || groupID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id and group_id are required"}}, nil
	}
	query := `mutation ($boardId: ID!, $groupId: String!) { delete_group(board_id: $boardId, group_id: $groupId) { id deleted } }`
	var result struct {
		DeleteGroup map[string]any `json:"delete_group"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardId": boardID, "groupId": groupID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.DeleteGroup}, nil
}
