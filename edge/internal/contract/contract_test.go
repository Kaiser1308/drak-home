package contract

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// exampleSchema maps every JSON example in contracts/v1/examples to the schema
// it must validate against. Examples missing from both this map and
// skippedExamples fail the test, so new examples cannot be silently ignored.
var exampleSchema = map[string]string{
	"artifact.json":           "artifact.json",
	"automation-status.json":  "automation-status.json",
	"automation.json":         "automation.json",
	"availability.json":       "availability.json",
	"command-result.json":     "command-result.json",
	"command.json":            "command.json",
	"deployment-commit.json":  "deployment.json",
	"deployment-prepare.json": "deployment.json",
	"deployment-result.json":  "deployment-result.json",
	"event.json":              "event.json",
	"health.json":             "health.json",
	"manifest.json":           "manifest.json",
	"state.json":              "state.json",
	"telemetry.json":          "telemetry.json",
}

// skippedExamples have no schema in contract v1 (reason given).
var skippedExamples = map[string]string{
	"ir.json":                       "compiler IR illustration; no v1 schema",
	"verification-certificate.json": "certificate illustration; no v1 schema",
	"automation.yaml":               "YAML source form of automation.json; S0 has no YAML parser",
}

func mustDir(t *testing.T) string {
	t.Helper()
	dir, err := Dir()
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func compileSchemas(t *testing.T, dir string) map[string]*jsonschema.Schema {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "schemas", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no schemas found")
	}
	c := jsonschema.NewCompiler()
	c.DefaultDraft(jsonschema.Draft2020)
	out := map[string]*jsonschema.Schema{}
	for _, f := range files {
		sch, err := c.Compile(f)
		if err != nil {
			t.Fatalf("compile %s: %v", filepath.Base(f), err)
		}
		out[filepath.Base(f)] = sch
	}
	return out
}

// TestAcceptanceContractExamplesValidate checks that every schema in
// contracts/v1/schemas compiles under Draft 2020-12 and every mapped example
// validates against its schema.
func TestAcceptanceContractExamplesValidate(t *testing.T) {
	dir := mustDir(t)
	schemas := compileSchemas(t, dir)

	for ex, sch := range exampleSchema {
		if _, ok := schemas[sch]; !ok {
			t.Errorf("example %s maps to missing schema %s", ex, sch)
		}
	}

	entries, err := os.ReadDir(filepath.Join(dir, "examples"))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	validated := 0
	for _, name := range names {
		if reason, skip := skippedExamples[name]; skip {
			t.Logf("skip %s: %s", name, reason)
			continue
		}
		schName, ok := exampleSchema[name]
		if !ok {
			t.Errorf("example %s has no schema mapping; add it to exampleSchema or skippedExamples", name)
			continue
		}
		t.Run(strings.TrimSuffix(name, ".json"), func(t *testing.T) {
			f, err := os.Open(filepath.Join(dir, "examples", name))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			doc, err := jsonschema.UnmarshalJSON(f)
			if err != nil {
				t.Fatal(err)
			}
			if err := schemas[schName].Validate(doc); err != nil {
				t.Errorf("%s does not validate against %s: %v", name, schName, err)
			}
		})
		validated++
	}
	if validated != len(exampleSchema) {
		t.Errorf("validated %d examples, want %d", validated, len(exampleSchema))
	}
}

// TestAcceptanceSchemaRejectsUnknownField proves validation is not vacuous:
// a command with an extra field must be rejected (additionalProperties=false).
func TestAcceptanceSchemaRejectsUnknownField(t *testing.T) {
	dir := mustDir(t)
	sch := compileSchemas(t, dir)["command.json"]
	doc, err := jsonschema.UnmarshalJSON(strings.NewReader(
		`{"command_id":"c","target_boot_id":"b","actuator":"fan","state":"on","ttl_ms":1,"extra":1}`))
	if err != nil {
		t.Fatal(err)
	}
	if err := sch.Validate(doc); err == nil {
		t.Fatal("command with unknown field validated; expected rejection")
	}
}

// TestLoadConformanceScenarios loads the shared test vectors fixture.
func TestLoadConformanceScenarios(t *testing.T) {
	s, err := LoadScenarios(mustDir(t))
	if err != nil {
		t.Fatal(err)
	}
	if s.Clock != "virtual_monotonic_ms" {
		t.Errorf("clock = %q, want virtual_monotonic_ms", s.Clock)
	}
	if len(s.Scenarios) == 0 {
		t.Fatal("no scenarios")
	}
	seen := map[string]bool{}
	for _, sc := range s.Scenarios {
		if sc.ID == "" || seen[sc.ID] {
			t.Errorf("scenario id %q empty or duplicated", sc.ID)
		}
		seen[sc.ID] = true
		if len(sc.Events) == 0 {
			t.Errorf("%s: no events", sc.ID)
		}
		for i := 1; i < len(sc.Events); i++ {
			if sc.Events[i].TMs < sc.Events[i-1].TMs {
				t.Errorf("%s: events not in time order at %d", sc.ID, i)
			}
		}
	}
	if _, ok := seen["T01"]; !ok {
		t.Error("scenario T01 missing")
	}
}

func TestDirEnvOverride(t *testing.T) {
	want := t.TempDir()
	t.Setenv(EnvDir, want)
	got, err := Dir()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("Dir() = %q, want %q", got, want)
	}
}
