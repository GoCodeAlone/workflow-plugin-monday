package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// Version is set at build time via -ldflags
// "-X github.com/GoCodeAlone/workflow-plugin-monday/internal.Version=X.Y.Z".
// Default is a bare semver so plugin loaders that validate semver accept
// unreleased dev builds; goreleaser overrides with the real release tag.
var Version = "0.0.0"

type mondayPlugin struct{}

func NewMondayPlugin() sdk.PluginProvider {
	return &mondayPlugin{}
}

func (p *mondayPlugin) Manifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Name:        "workflow-plugin-monday",
		Version:     Version,
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
