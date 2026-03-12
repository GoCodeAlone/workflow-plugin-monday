package internal

import (
	"context"
	"testing"
)

func TestMondayModule_MissingAPIToken(t *testing.T) {
	_, err := newMondayModule("test", map[string]any{})
	if err == nil {
		t.Error("expected error for missing apiToken")
	}
}

func TestMondayModule_InitRegistersClient(t *testing.T) {
	m, err := newMondayModule("test-init", map[string]any{"apiToken": "tok_test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	defer m.Stop(context.Background()) //nolint:errcheck

	if _, ok := GetClient("test-init"); !ok {
		t.Error("expected client to be registered after Init")
	}
}

func TestMondayModule_StopUnregisters(t *testing.T) {
	RegisterClient("test-stop", NewMondayClient("tok"))

	m := &mondayModule{name: "test-stop", apiToken: "tok"}
	if err := m.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := GetClient("test-stop"); ok {
		t.Error("expected client to be unregistered after Stop")
	}
}

func TestMondayModule_StartIsNoop(t *testing.T) {
	m := &mondayModule{name: "test", apiToken: "tok"}
	if err := m.Start(context.Background()); err != nil {
		t.Errorf("unexpected error from Start: %v", err)
	}
}

func TestMondayModule_WithAPIVersion(t *testing.T) {
	m, err := newMondayModule("test-ver", map[string]any{
		"apiToken":   "tok",
		"apiVersion": "2025-01",
	})
	if err != nil {
		t.Fatal(err)
	}
	if m.apiVersion != "2025-01" {
		t.Errorf("expected apiVersion 2025-01, got %s", m.apiVersion)
	}
}
