package internal

import (
	"context"
	"fmt"
)

type mondayModule struct {
	name       string
	apiToken   string
	apiVersion string
	baseURL    string
}

func newMondayModule(name string, config map[string]any) (*mondayModule, error) {
	apiToken, _ := config["apiToken"].(string)
	if apiToken == "" {
		return nil, fmt.Errorf("monday.provider %q: apiToken is required", name)
	}
	m := &mondayModule{
		name:     name,
		apiToken: apiToken,
	}
	if v, ok := config["apiVersion"].(string); ok && v != "" {
		m.apiVersion = v
	}
	if v, ok := config["baseUrl"].(string); ok && v != "" {
		m.baseURL = v
	}
	return m, nil
}

func (m *mondayModule) Init() error {
	opts := []ClientOption{}
	if m.apiVersion != "" {
		opts = append(opts, WithAPIVersion(m.apiVersion))
	}
	if m.baseURL != "" {
		opts = append(opts, WithBaseURL(m.baseURL))
	}
	client := NewMondayClient(m.apiToken, opts...)
	RegisterClient(m.name, client)
	return nil
}

func (m *mondayModule) Start(_ context.Context) error { return nil }

func (m *mondayModule) Stop(_ context.Context) error {
	UnregisterClient(m.name)
	return nil
}
