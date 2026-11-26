// Package contract locates and loads the shared DRAK contract v1 files
// (schemas, examples, conformance scenarios) from the repository's
// contracts/v1 directory. Edge and device read the same files; nothing is
// vendored or pinned by hash.
package contract

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// relDir is the contract directory relative to the repository root.
const relDir = "contracts/v1"

// EnvDir overrides contract directory discovery when set.
const EnvDir = "DRAK_CONTRACTS_DIR"

// Dir returns the absolute path of contracts/v1. It honours DRAK_CONTRACTS_DIR,
// otherwise walks up from the working directory until it finds contracts/v1.
func Dir() (string, error) {
	if d := os.Getenv(EnvDir); d != "" {
		return filepath.Abs(d)
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for dir := wd; ; {
		cand := filepath.Join(dir, relDir)
		if st, err := os.Stat(filepath.Join(cand, "schemas")); err == nil && st.IsDir() {
			return cand, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("contract: contracts/v1 not found above " + wd)
		}
		dir = parent
	}
}

// Event is one input of a conformance scenario.
type Event struct {
	TMs               int64  `json:"t_ms"`
	Type              string `json:"type"`
	ValueMilliCelsius *int64 `json:"value_milli_celsius,omitempty"`
	Enabled           *bool  `json:"enabled,omitempty"`
	Online            *bool  `json:"online,omitempty"`
}

// Action is an expected actuator write of a conformance scenario.
type Action struct {
	TMs      int64  `json:"t_ms"`
	Actuator string `json:"actuator"`
	State    string `json:"state"`
}

// Scenario is one conformance test vector.
type Scenario struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	InitialFan  string `json:"initial_fan"`
	// ArtifactOverrides patches fields of the referenced artifact for this
	// scenario only; shape is left open by the contract.
	ArtifactOverrides map[string]any `json:"artifact_overrides,omitempty"`
	Events            []Event        `json:"events"`
	ExpectedActions   []Action       `json:"expected_actions"`
}

// Scenarios is conformance/scenarios.json.
type Scenarios struct {
	Status      string     `json:"status"`
	Artifact    string     `json:"artifact"`
	Clock       string     `json:"clock"`
	Assumptions []string   `json:"assumptions"`
	Scenarios   []Scenario `json:"scenarios"`
}

// LoadScenarios reads conformance/scenarios.json from dir (a contracts/v1 path).
func LoadScenarios(dir string) (*Scenarios, error) {
	b, err := os.ReadFile(filepath.Join(dir, "conformance", "scenarios.json"))
	if err != nil {
		return nil, err
	}
	var s Scenarios
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&s); err != nil {
		return nil, fmt.Errorf("contract: scenarios.json: %w", err)
	}
	return &s, nil
}
