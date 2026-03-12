package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type mondayPlugin struct{}

func NewMondayPlugin() sdk.PluginProvider {
	return &mondayPlugin{}
}

func (p *mondayPlugin) Manifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Name:        "workflow-plugin-monday",
		Version:     "0.1.0",
		Author:      "GoCodeAlone",
		Description: "monday.com integration plugin (~57 step types covering all monday.com resources)",
	}
}

func (p *mondayPlugin) ModuleTypes() []string {
	return []string{"monday.provider"}
}

func (p *mondayPlugin) CreateModule(typeName, name string, config map[string]any) (sdk.ModuleInstance, error) {
	switch typeName {
	case "monday.provider":
		m, err := newMondayModule(name, config)
		if err != nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, fmt.Errorf("monday plugin: unknown module type %q", typeName)
	}
}

func (p *mondayPlugin) StepTypes() []string {
	return allStepTypes()
}

func (p *mondayPlugin) CreateStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	constructor, ok := stepRegistry()[typeName]
	if !ok {
		return nil, fmt.Errorf("monday plugin: unknown step type %q", typeName)
	}
	return constructor(name, config)
}
