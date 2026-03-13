package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type createWebhookStep struct{ name, moduleName string }

func newCreateWebhookStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createWebhookStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createWebhookStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	url := resolveValue("url", current, config)
	event := resolveValue("event", current, config)
	if boardID == "" || url == "" || event == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id, url, and event are required"}}, nil
	}
	query := `mutation ($boardId: ID!, $url: String!, $event: WebhookEventType!) {
		create_webhook(board_id: $boardId, url: $url, event: $event) { id board_id event }
	}`
	var result struct {
		CreateWebhook map[string]any `json:"create_webhook"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardId": boardID, "url": url, "event": event}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateWebhook}, nil
}

type listWebhooksStep struct{ name, moduleName string }

func newListWebhooksStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listWebhooksStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listWebhooksStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	boardID := resolveValue("board_id", current, config)
	if boardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "board_id is required"}}, nil
	}
	query := `query ($boardId: ID!) { webhooks(board_id: $boardId) { id board_id event url } }`
	var result struct {
		Webhooks []map[string]any `json:"webhooks"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"boardId": boardID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"webhooks": toAnySlice(result.Webhooks)}}, nil
}

type deleteWebhookStep struct{ name, moduleName string }

func newDeleteWebhookStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &deleteWebhookStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteWebhookStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	webhookID := resolveValue("webhook_id", current, config)
	if webhookID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "webhook_id is required"}}, nil
	}
	query := `mutation ($id: ID!) { delete_webhook(id: $id) { id board_id } }`
	var result struct {
		DeleteWebhook map[string]any `json:"delete_webhook"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"id": webhookID}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.DeleteWebhook}, nil
}
