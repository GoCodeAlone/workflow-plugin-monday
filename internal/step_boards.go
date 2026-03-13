package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// createBoardStep implements step.monday_create_board
type createBoardStep struct{ name, moduleName string }

func newCreateBoardStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createBoardStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createBoardStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardName := resolveValue("board_name", current, config)
	if boardName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_name is required"}}, nil
	}
	boardKind := resolveValue("board_kind", current, config)
	if boardKind == "" {
		boardKind = "public"
	}
	query := `mutation ($boardName: String!, $boardKind: BoardKind!) {
		create_board(board_name: $boardName, board_kind: $boardKind) { id name board_kind state }
	}`
	var result struct {
		CreateBoard map[string]any `json:"create_board"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardName": boardName, "boardKind": boardKind}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateBoard}, nil
}

// listBoardsStep implements step.monday_list_boards
type listBoardsStep struct{ name, moduleName string }

func newListBoardsStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listBoardsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listBoardsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	limit := resolveInt("limit", current, config)
	if limit == 0 {
		limit = 25
	}
	query := `query ($limit: Int!) { boards(limit: $limit) { id name board_kind state description } }`
	var result struct {
		Boards []any `json:"boards"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"limit": limit}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"boards": result.Boards}}, nil
}

// fetchBoardStep implements step.monday_fetch_board
type fetchBoardStep struct{ name, moduleName string }

func newFetchBoardStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &fetchBoardStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *fetchBoardStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	if boardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id is required"}}, nil
	}
	query := `query ($ids: [ID!]) { boards(ids: $ids) { id name board_kind state description } }`
	var result struct {
		Boards []any `json:"boards"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"ids": []string{boardID}}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if len(result.Boards) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "board not found"}}, nil
	}
	return &sdk.StepResult{Output: result.Boards[0].(map[string]any)}, nil
}

// updateBoardStep (fetchBoardStep returns single map - no slice fix needed) implements step.monday_update_board
type updateBoardStep struct{ name, moduleName string }

func newUpdateBoardStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &updateBoardStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *updateBoardStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	attribute := resolveValue("attribute", current, config)
	newValue := resolveValue("new_value", current, config)
	if boardID == "" || attribute == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id and attribute are required"}}, nil
	}
	query := `mutation ($boardId: ID!, $boardAttribute: BoardAttributes!, $newValue: String!) {
		update_board(board_id: $boardId, board_attribute: $boardAttribute, new_value: $newValue)
	}`
	data, err := client.Execute(ctx, query, map[string]any{"boardId": boardID, "boardAttribute": attribute, "newValue": newValue})
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": string(data)}}, nil
}

// deleteBoardStep implements step.monday_delete_board
type deleteBoardStep struct{ name, moduleName string }

func newDeleteBoardStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &deleteBoardStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteBoardStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	if boardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id is required"}}, nil
	}
	query := `mutation ($boardId: ID!) { delete_board(board_id: $boardId) { id } }`
	var result struct {
		DeleteBoard map[string]any `json:"delete_board"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardId": boardID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.DeleteBoard}, nil
}

// duplicateBoardStep implements step.monday_duplicate_board
type duplicateBoardStep struct{ name, moduleName string }

func newDuplicateBoardStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &duplicateBoardStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *duplicateBoardStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	if boardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id is required"}}, nil
	}
	duplicateType := resolveValue("duplicate_type", current, config)
	if duplicateType == "" {
		duplicateType = "duplicate_board_with_structure"
	}
	query := `mutation ($boardId: ID!, $duplicateType: DuplicateBoardType!) {
		duplicate_board(board_id: $boardId, duplicate_type: $duplicateType) { board { id name } }
	}`
	var result struct {
		DuplicateBoard struct {
			Board map[string]any `json:"board"`
		} `json:"duplicate_board"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardId": boardID, "duplicateType": duplicateType}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.DuplicateBoard.Board}, nil
}

// archiveBoardStep implements step.monday_archive_board
type archiveBoardStep struct{ name, moduleName string }

func newArchiveBoardStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &archiveBoardStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *archiveBoardStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	if boardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id is required"}}, nil
	}
	query := `mutation ($boardId: ID!) { archive_board(board_id: $boardId) { id state } }`
	var result struct {
		ArchiveBoard map[string]any `json:"archive_board"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardId": boardID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.ArchiveBoard}, nil
}
