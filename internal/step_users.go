package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type listUsersStep struct{ name, moduleName string }

func newListUsersStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listUsersStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listUsersStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	limit := resolveInt("limit", current, config)
	if limit == 0 {
		limit = 25
	}
	query := `query ($limit: Int!) { users(limit: $limit) { id name email } }`
	var result struct {
		Users []map[string]any `json:"users"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"limit": limit}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"users": toAnySlice(result.Users)}}, nil
}

type fetchUserStep struct{ name, moduleName string }

func newFetchUserStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &fetchUserStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *fetchUserStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	userID := resolveValue("user_id", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "user_id is required"}}, nil
	}
	query := `query ($ids: [ID!]) { users(ids: $ids) { id name email } }`
	var result struct {
		Users []map[string]any `json:"users"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"ids": []string{userID}}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if len(result.Users) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "user not found"}}, nil
	}
	return &sdk.StepResult{Output: result.Users[0]}, nil
}

type inviteUserStep struct{ name, moduleName string }

func newInviteUserStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &inviteUserStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *inviteUserStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	email := resolveValue("email", current, config)
	userKind := resolveValue("user_kind", current, config)
	if email == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "email is required"}}, nil
	}
	if userKind == "" {
		userKind = "member"
	}
	query := `mutation ($email: String!, $userKind: UserKind!) {
		invite_user_to_account(email: $email, user_kind: $userKind) { id name email }
	}`
	var result struct {
		InviteUser map[string]any `json:"invite_user_to_account"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"email": email, "userKind": userKind}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.InviteUser}, nil
}
