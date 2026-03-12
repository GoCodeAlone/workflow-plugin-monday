package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type createWorkspaceStep struct{ name, moduleName string }

func newCreateWorkspaceStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createWorkspaceStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createWorkspaceStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	wsName := resolveValue("name", current, config)
	wsKind := resolveValue("kind", current, config)
	if wsName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	if wsKind == "" {
		wsKind = "open"
	}
	query := `mutation ($name: String!, $kind: WorkspaceKind!) {
		create_workspace(name: $name, kind: $kind) { id name kind }
	}`
	var result struct {
		CreateWorkspace map[string]any `json:"create_workspace"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"name": wsName, "kind": wsKind}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateWorkspace}, nil
}

type listWorkspacesStep struct{ name, moduleName string }

func newListWorkspacesStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listWorkspacesStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listWorkspacesStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	limit := resolveInt("limit", current, config)
	if limit == 0 {
		limit = 25
	}
	query := `query ($limit: Int!) { workspaces(limit: $limit) { id name kind } }`
	var result struct {
		Workspaces []map[string]any `json:"workspaces"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"limit": limit}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"workspaces": result.Workspaces}}, nil
}

type updateWorkspaceStep struct{ name, moduleName string }

func newUpdateWorkspaceStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &updateWorkspaceStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *updateWorkspaceStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	wsID := resolveValue("workspace_id", current, config)
	attribute := resolveValue("attribute", current, config)
	newValue := resolveValue("new_value", current, config)
	if wsID == "" || attribute == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "workspace_id and attribute are required"}}, nil
	}
	query := `mutation ($id: ID!, $attribute: WorkspaceAttributes!, $value: String!) {
		update_workspace(id: $id, attributes: {name: $value}) { id name }
	}`
	var result struct {
		UpdateWorkspace map[string]any `json:"update_workspace"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"id": wsID, "attribute": attribute, "value": newValue}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.UpdateWorkspace}, nil
}

type deleteWorkspaceStep struct{ name, moduleName string }

func newDeleteWorkspaceStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &deleteWorkspaceStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteWorkspaceStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	wsID := resolveValue("workspace_id", current, config)
	if wsID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "workspace_id is required"}}, nil
	}
	query := `mutation ($id: ID!) { delete_workspace(workspace_id: $id) { id } }`
	var result struct {
		DeleteWorkspace map[string]any `json:"delete_workspace"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"id": wsID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.DeleteWorkspace}, nil
}
