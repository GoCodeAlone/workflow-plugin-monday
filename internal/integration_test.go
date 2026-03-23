package internal_test

import (
	"testing"

	"github.com/GoCodeAlone/workflow/wftest"
)

func TestIntegration_CreateItem(t *testing.T) {
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  create:
    steps:
      - name: item
        type: step.monday_create_item
        config:
          board_id: "123456"
          item_name: "New Task"
      - name: confirm
        type: step.set
        config:
          values:
            created: true
`),
		wftest.MockStep("step.monday_create_item", wftest.Returns(map[string]any{
			"id": "987654", "name": "New Task",
		})),
	)
	result := h.ExecutePipeline("create", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["created"] != true {
		t.Errorf("expected created=true, got %v", result.Output["created"])
	}
	itemStep := result.StepResults["item"]
	if itemStep == nil {
		t.Fatal("expected step output for 'item'")
	}
	if itemStep["id"] != "987654" {
		t.Errorf("expected item id=987654, got %v", itemStep["id"])
	}
}

func TestIntegration_ListItems(t *testing.T) {
	items := []any{
		map[string]any{"id": "1", "name": "Task A", "state": "active"},
		map[string]any{"id": "2", "name": "Task B", "state": "active"},
	}
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  list:
    steps:
      - name: items
        type: step.monday_list_items
        config:
          board_id: "123456"
          limit: 10
      - name: confirm
        type: step.set
        config:
          values:
            listed: true
`),
		wftest.MockStep("step.monday_list_items", wftest.Returns(map[string]any{
			"items": items,
		})),
	)
	result := h.ExecutePipeline("list", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["listed"] != true {
		t.Errorf("expected listed=true, got %v", result.Output["listed"])
	}
	itemsStep := result.StepResults["items"]
	if itemsStep == nil {
		t.Fatal("expected step output for 'items'")
	}
	got, ok := itemsStep["items"].([]any)
	if !ok {
		t.Fatalf("expected items slice in step output, got %T", itemsStep["items"])
	}
	if len(got) != 2 {
		t.Errorf("expected 2 items, got %d", len(got))
	}
}

func TestIntegration_CreateBoard(t *testing.T) {
	rec := wftest.RecordStep("step.monday_create_board")
	rec.WithOutput(map[string]any{"id": "555", "name": "My Board", "state": "active"})

	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  setup:
    steps:
      - name: board
        type: step.monday_create_board
        config:
          board_name: "My Board"
          board_kind: "public"
      - name: confirm
        type: step.set
        config:
          values:
            setup_done: true
`),
		rec,
	)
	result := h.ExecutePipeline("setup", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["setup_done"] != true {
		t.Errorf("expected setup_done=true, got %v", result.Output["setup_done"])
	}
	if rec.CallCount() != 1 {
		t.Errorf("expected step.monday_create_board to be called once, got %d", rec.CallCount())
	}
	boardStep := result.StepResults["board"]
	if boardStep == nil {
		t.Fatal("expected step output for 'board'")
	}
	if boardStep["id"] != "555" {
		t.Errorf("expected board id=555, got %v", boardStep["id"])
	}
}
