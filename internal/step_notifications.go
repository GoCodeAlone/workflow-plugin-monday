package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type createNotificationStep struct{ name, moduleName string }

func newCreateNotificationStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createNotificationStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createNotificationStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	userID := resolveValue("user_id", current, config)
	targetID := resolveValue("target_id", current, config)
	text := resolveValue("text", current, config)
	targetType := resolveValue("target_type", current, config)
	if userID == "" || targetID == "" || text == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "user_id, target_id, and text are required"}}, nil
	}
	if targetType == "" {
		targetType = "Project"
	}
	query := `mutation ($userId: ID!, $targetId: ID!, $text: String!, $targetType: NotificationTargetType!) {
		create_notification(user_id: $userId, target_id: $targetId, text: $text, target_type: $targetType) { id text }
	}`
	var result struct {
		CreateNotification map[string]any `json:"create_notification"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{
		"userId":     userID,
		"targetId":   targetID,
		"text":       text,
		"targetType": targetType,
	}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateNotification}, nil
}
