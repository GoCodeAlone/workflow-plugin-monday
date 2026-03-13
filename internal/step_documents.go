package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type createDocumentStep struct{ name, moduleName string }

func newCreateDocumentStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &createDocumentStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createDocumentStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	location := resolveValue("location", current, config)
	title := resolveValue("title", current, config)
	if location == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "location is required"}}, nil
	}
	if title == "" {
		title = "New Document"
	}
	query := `mutation ($location: DocumentLocation!, $title: String!) {
		create_doc(location: {workspace: {workspace_id: $location}}, title: $title) { id object_id }
	}`
	var result struct {
		CreateDoc map[string]any `json:"create_doc"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"location": location, "title": title}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.CreateDoc}, nil
}

type listDocumentsStep struct{ name, moduleName string }

func newListDocumentsStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &listDocumentsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listDocumentsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	limit := resolveInt("limit", current, config)
	if limit == 0 {
		limit = 25
	}
	query := `query ($limit: Int!) { docs(limit: $limit) { id object_id title } }`
	var result struct {
		Docs []map[string]any `json:"docs"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"limit": limit}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"documents": toAnySlice(result.Docs)}}, nil
}

type updateDocumentStep struct{ name, moduleName string }

func newUpdateDocumentStep(name string, config map[string]any) (sdk.StepInstance, error) {
	return &updateDocumentStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *updateDocumentStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "monday client not found: " + s.moduleName}}, nil
	}
	docID := resolveValue("doc_id", current, config)
	content := resolveValue("content", current, config)
	if docID == "" || content == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "doc_id and content are required"}}, nil
	}
	query := `mutation ($docId: ID!, $content: JSON!) {
		add_doc_block(doc_id: $docId, type: normal_text, content: $content) { id type }
	}`
	var result struct {
		AddDocBlock map[string]any `json:"add_doc_block"`
	}
	if err := client.ExecuteInto(ctx, query, map[string]any{"docId": docID, "content": content}, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result.AddDocBlock}, nil
}
