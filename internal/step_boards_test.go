package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeTestServer(t *testing.T, response map[string]any) (*httptest.Server, *MondayClient) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response) //nolint:errcheck
	}))
	t.Cleanup(srv.Close)
	client := NewMondayClient("test-token", WithBaseURL(srv.URL))
	return srv, client
}

func makeErrorServer(t *testing.T, gqlError string) (*httptest.Server, *MondayClient) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"errors": []map[string]any{{"message": gqlError}},
		})
	}))
	t.Cleanup(srv.Close)
	client := NewMondayClient("test-token", WithBaseURL(srv.URL))
	return srv, client
}

func TestCreateBoardStep_MissingBoardName(t *testing.T) {
	step := &createBoardStep{name: "test", moduleName: "test-mod"}
	result, err := step.Execute(t.Context(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing board_name")
	}
}

func TestCreateBoardStep_MissingClient(t *testing.T) {
	step := &createBoardStep{name: "test", moduleName: "nonexistent"}
	result, err := step.Execute(t.Context(), nil, nil, map[string]any{"board_name": "Test"}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing client")
	}
}

func TestCreateBoardStep_GraphQLError(t *testing.T) {
	_, client := makeErrorServer(t, "not authorized")
	RegisterClient("test-boards-err", client)
	defer UnregisterClient("test-boards-err")

	step := &createBoardStep{name: "test", moduleName: "test-boards-err"}
	result, err := step.Execute(t.Context(), nil, nil, map[string]any{"board_name": "Test"}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error output for GraphQL error")
	}
}

func TestCreateBoardStep_Success(t *testing.T) {
	_, client := makeTestServer(t, map[string]any{
		"data": map[string]any{
			"create_board": map[string]any{"id": "123", "name": "My Board", "board_kind": "public", "state": "active"},
		},
	})
	RegisterClient("test-boards-ok", client)
	defer UnregisterClient("test-boards-ok")

	step := &createBoardStep{name: "test", moduleName: "test-boards-ok"}
	result, err := step.Execute(t.Context(), nil, nil, map[string]any{"board_name": "My Board"}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] != nil {
		t.Errorf("unexpected error: %v", result.Output["error"])
	}
	if result.Output["id"] != "123" {
		t.Errorf("expected id=123, got %v", result.Output["id"])
	}
}

func TestListBoardsStep_Success(t *testing.T) {
	_, client := makeTestServer(t, map[string]any{
		"data": map[string]any{
			"boards": []map[string]any{
				{"id": "1", "name": "Board A"},
				{"id": "2", "name": "Board B"},
			},
		},
	})
	RegisterClient("test-list-boards", client)
	defer UnregisterClient("test-list-boards")

	step := &listBoardsStep{name: "test", moduleName: "test-list-boards"}
	result, err := step.Execute(t.Context(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	boards, ok := result.Output["boards"].([]any)
	if !ok {
		t.Fatalf("expected boards slice, got %T", result.Output["boards"])
	}
	if len(boards) != 2 {
		t.Errorf("expected 2 boards, got %d", len(boards))
	}
}

func TestDeleteBoardStep_MissingBoardID(t *testing.T) {
	step := &deleteBoardStep{name: "test", moduleName: "x"}
	result, err := step.Execute(t.Context(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing board_id")
	}
}
