package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type listTeamsStep struct{ name, moduleName string }

func newListTeamsStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listTeamsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listTeamsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, _ map[string]any, _ map[string]any, _ map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	query := `query { teams { id name } }`
	var result struct {
		Teams []map[string]any `json:"teams"`
	}
	if err := client.ExecuteInto(ctx, query, nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"teams": result.Teams}}, nil
}

type addTeamToWorkspaceStep struct{ name, moduleName string }

func newAddTeamToWorkspaceStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &addTeamToWorkspaceStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *addTeamToWorkspaceStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	wsID := resolveValue("workspace_id", current, config)
	teamID := resolveValue("team_id", current, config)
	if wsID == "" || teamID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "workspace_id and team_id are required"}}, nil
	}
	query := `mutation ($workspaceId: ID!, $teamId: ID!) {
		add_teams_to_workspace(workspace_id: $workspaceId, team_ids: [$teamId]) { id name }
	}`
	var result struct {
		AddTeams map[string]any `json:"add_teams_to_workspace"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"workspaceId": wsID, "teamId": teamID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.AddTeams}, nil
}
