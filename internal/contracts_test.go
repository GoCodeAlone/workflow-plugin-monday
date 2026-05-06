package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type manifestFile struct {
	Version string `json:"version"`
	Capabilities struct {
		ModuleTypes []string `json:"moduleTypes"`
		StepTypes   []string `json:"stepTypes"`
	} `json:"capabilities"`
	Downloads []struct {
		OS   string `json:"os"`
		Arch string `json:"arch"`
		URL  string `json:"url"`
	} `json:"downloads"`
}

type contractsFile struct {
	Contracts []contractDescriptor `json:"contracts"`
}

type contractDescriptor struct {
	Kind string `json:"kind"`
	Type string `json:"type"`
	Mode string `json:"mode"`
}

func TestPluginManifestAndContractsMatchRuntimeTypes(t *testing.T) {
	manifest := readJSONFile[manifestFile](t, "../plugin.json")
	contracts := readJSONFile[contractsFile](t, "../plugin.contracts.json")
	plugin := &mondayPlugin{}

	assertStringSetEqual(t, "manifest moduleTypes", manifest.Capabilities.ModuleTypes, plugin.ModuleTypes())
	assertStringSetEqual(t, "manifest stepTypes", manifest.Capabilities.StepTypes, plugin.StepTypes())
	assertStringSetEqual(t, "module contracts", contractTypes(t, contracts.Contracts, "module"), plugin.ModuleTypes())
	assertStringSetEqual(t, "step contracts", contractTypes(t, contracts.Contracts, "step"), plugin.StepTypes())
}

func TestPluginManifestVersionAndDownloadsConsistent(t *testing.T) {
	manifest := readJSONFile[manifestFile](t, "../plugin.json")

	if manifest.Version == "" {
		t.Fatal("plugin.json: version field is empty")
	}

	if len(manifest.Downloads) == 0 {
		t.Fatal("plugin.json: downloads field is empty")
	}

	wantTag := fmt.Sprintf("/v%s/", manifest.Version)
	for _, dl := range manifest.Downloads {
		if !strings.Contains(dl.URL, wantTag) {
			t.Errorf("plugin.json: download URL %q does not contain version tag %s", dl.URL, wantTag)
		}
	}
}

func readJSONFile[T any](t *testing.T, path string) T {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return out
}

func contractTypes(t *testing.T, contracts []contractDescriptor, kind string) []string {
	t.Helper()

	var types []string
	for _, contract := range contracts {
		if contract.Kind != kind {
			continue
		}
		if contract.Mode != "strict" {
			t.Fatalf("%s contract %q uses mode %q, want strict", kind, contract.Type, contract.Mode)
		}
		types = append(types, contract.Type)
	}
	return types
}

func assertStringSetEqual(t *testing.T, label string, got, want []string) {
	t.Helper()

	got = slices.Clone(got)
	want = slices.Clone(want)
	slices.Sort(got)
	slices.Sort(want)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s mismatch\ngot:  %v\nwant: %v", label, got, want)
	}
}
