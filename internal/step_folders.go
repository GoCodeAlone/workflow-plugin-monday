package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type createFolderStep struct{ name, moduleName string }

func newCreateFolderStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createFolderStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createFolderStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	wsID := resolveValue("workspace_id", current, config)
	folderName := resolveValue("name", current, config)
	if wsID == "" || folderName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "workspace_id and name are required"}}, nil
	}
	query := `mutation ($workspaceId: ID!, $name: String!) {
		create_folder(workspace_id: $workspaceId, name: $name) { id name }
	}`
	var result struct {
		CreateFolder map[string]any `json:"create_folder"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"workspaceId": wsID, "name": folderName}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateFolder}, nil
}

type listFoldersStep struct{ name, moduleName string }

func newListFoldersStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listFoldersStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listFoldersStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	wsID := resolveValue("workspace_id", current, config)
	vars := map[string]any{}
	if wsID != "" {
		vars["workspaceIds"] = []string{wsID}
	}
	query := `query ($workspaceIds: [ID]) { folders(workspace_ids: $workspaceIds) { id name } }`
	var result struct {
		Folders []map[string]any `json:"folders"`
	}
	if err := client.ExecuteInto(ctx, query, vars, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"folders": toAnySlice(result.Folders)}}, nil
}

type updateFolderStep struct{ name, moduleName string }

func newUpdateFolderStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &updateFolderStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *updateFolderStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	folderID := resolveValue("folder_id", current, config)
	newName := resolveValue("name", current, config)
	if folderID == "" || newName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "folder_id and name are required"}}, nil
	}
	query := `mutation ($folderId: ID!, $name: String!) {
		update_folder(folder_id: $folderId, name: $name) { id name }
	}`
	var result struct {
		UpdateFolder map[string]any `json:"update_folder"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"folderId": folderID, "name": newName}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.UpdateFolder}, nil
}

type deleteFolderStep struct{ name, moduleName string }

func newDeleteFolderStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &deleteFolderStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteFolderStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	folderID := resolveValue("folder_id", current, config)
	if folderID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "folder_id is required"}}, nil
	}
	query := `mutation ($folderId: ID!) { delete_folder(folder_id: $folderId) { id } }`
	var result struct {
		DeleteFolder map[string]any `json:"delete_folder"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"folderId": folderID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.DeleteFolder}, nil
}
