package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type uploadFileStep struct{ name, moduleName string }

func newUploadFileStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &uploadFileStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *uploadFileStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	columnID := resolveValue("column_id", current, config)
	fileURL := resolveValue("file_url", current, config)
	if itemID == "" || columnID == "" || fileURL == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id, column_id, and file_url are required"}}, nil
	}
	// monday.com file upload via GraphQL requires multipart — use add_file_to_column
	query := `mutation ($itemId: ID!, $columnId: String!, $fileUrl: String!) {
		add_file_to_column(item_id: $itemId, column_id: $columnId, file: $fileUrl) { id }
	}`
	var result struct {
		AddFileToColumn map[string]any `json:"add_file_to_column"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"itemId": itemID, "columnId": columnID, "fileUrl": fileURL}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.AddFileToColumn}, nil
}

type listFilesStep struct{ name, moduleName string }

func newListFilesStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listFilesStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listFilesStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	itemID := resolveValue("item_id", current, config)
	if itemID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "item_id is required"}}, nil
	}
	limit := resolveInt("limit", current, config)
	if limit == 0 {
		limit = 25
	}
	query := `query ($ids: [ID!], $limit: Int!) {
		items(ids: $ids) { assets(limit: $limit) { id name url } }
	}`
	var result struct {
		Items []struct {
			Assets []map[string]any `json:"assets"`
		} `json:"items"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"ids": []string{itemID}, "limit": limit}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if len(result.Items) == 0 {
		return &sdk.StepResult{Output: map[string]any{"files": []any{}}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"files": result.Items[0].Assets}}, nil
}
