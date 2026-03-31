package configuration

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/zukrin/versifyr/internal/logging"
)

func TestGetBasePath(t *testing.T) {
	// Test default
	os.Unsetenv(BASEPATH_ENV)
	if bp := GetBasePath(); bp != BASEPATH_DEFAULT {
		t.Errorf("expected default %s, got %s", BASEPATH_DEFAULT, bp)
	}

	// Test env override
	expected := "custom-path"
	os.Setenv(BASEPATH_ENV, expected)
	defer os.Unsetenv(BASEPATH_ENV)
	if bp := GetBasePath(); bp != expected {
		t.Errorf("expected custom %s, got %s", expected, bp)
	}
}

func TestNewConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "versifyr-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	confDir := filepath.Join(tmpDir, BASEPATH_DEFAULT)
	if err := os.Mkdir(confDir, 0755); err != nil {
		t.Fatal(err)
	}

	confFile := filepath.Join(confDir, CONFIG_FILENAME)
	content := `
files:
  - name: test.go
    type: go
    path: test.go
`
	if err := os.WriteFile(confFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Change working directory to tmpDir to test NewConfig
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	cfg := &Config{}
	err = NewConfig(cfg)
	if err != nil {
		t.Fatalf("NewConfig failed: %v", err)
	}

	if len(cfg.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(cfg.Files))
	}
	if cfg.Files[0].Name != "test.go" {
		t.Errorf("expected test.go, got %s", cfg.Files[0].Name)
	}
}

func TestCompilePatternsEmbedded(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "versifyr-compile-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFilePath := filepath.Join(tmpDir, "version.go")
	testFileContent := `package main
// $versifyr:template={{ .version }}$
const Version = "v0.0.1"
`
	if err := os.WriteFile(testFilePath, []byte(testFileContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Files: []*ConfigFile{
			{
				Name: "version.go",
				Path: testFilePath,
				Type: "go",
			},
		},
	}

	logger := logging.NewLogger()
	_, err = cfg.CompilePatterns(logger)
	if err != nil {
		t.Fatalf("CompilePatterns failed: %v", err)
	}

	if len(cfg.Files[0].Placeholders) != 1 {
		t.Errorf("expected 1 placeholder, got %d", len(cfg.Files[0].Placeholders))
	}

	ph := cfg.Files[0].Placeholders[0]
	if ph.TemplateText != "{{ .version }}" {
		t.Errorf("expected template {{ .version }}, got %s", ph.TemplateText)
	}
	if ph.Line != 2 { // line indices start at 0, template is on line 1, target is line 2
		t.Errorf("expected target line 2, got %d", ph.Line)
	}
}

func TestCompilePatternsExplicit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "versifyr-compile-explicit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFilePath := filepath.Join(tmpDir, "version.go")
	testFileContent := `package main
const Version = "v0.0.1"
`
	if err := os.WriteFile(testFilePath, []byte(testFileContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Files: []*ConfigFile{
			{
				Name: "version.go",
				Path: testFilePath,
				Type: "go",
				Templates: []*Template{
					{
						Row:      2,
						Template: "{{ .version }}",
					},
				},
			},
		},
	}

	logger := logging.NewLogger()
	_, err = cfg.CompilePatterns(logger)
	if err != nil {
		t.Fatalf("CompilePatterns failed: %v", err)
	}

	if len(cfg.Files[0].Placeholders) != 1 {
		t.Errorf("expected 1 placeholder, got %d", len(cfg.Files[0].Placeholders))
	}

	ph := cfg.Files[0].Placeholders[0]
	if ph.TemplateText != "{{ .version }}" {
		t.Errorf("expected template {{ .version }}, got %s", ph.TemplateText)
	}
	if ph.Line != 1 { // Row 2 corresponds to index 1
		t.Errorf("expected target line index 1, got %d", ph.Line)
	}
}

func TestApplyTemplates(t *testing.T) {
	cfg := &ConfigFile{
		Name: "test.go",
		Lines: []string{
			"package main",
			"const Version = \"v0.0.1\"",
			"const Tag = \"\"",
		},
		Unescape: false,
	}

	// Mocking patterns compilation manually for the test
	testTpl1 := "{{ .version }}"
	testTpl2 := "{{ .version | replace \".\" \"_\" }}"

	tpl1, _ := cfg.CompileTemplate(testTpl1)
	tpl2, _ := cfg.CompileTemplate(testTpl2)

	cfg.Placeholders = []*Placeholder{
		{TemplateText: testTpl1, Template: tpl1, Line: 1},
		{TemplateText: testTpl2, Template: tpl2, Line: 2},
	}

	dictionary := map[string]string{
		"version": "v1.2.3",
	}

	err := cfg.ApplyTemplates(dictionary)
	if err != nil {
		t.Fatalf("ApplyTemplates failed: %v", err)
	}

	if cfg.Lines[1] != "v1.2.3" {
		t.Errorf("expected v1.2.3, got %s", cfg.Lines[1])
	}
	if cfg.Lines[2] != "v1_2_3" {
		t.Errorf("expected v1_2_3, got %s", cfg.Lines[2])
	}

	// Test Unescape
	cfgUnescape := &ConfigFile{
		Name:     "test.json",
		Lines:    []string{"{\"version\": \"\"}"},
		Unescape: true,
	}
	tplJson, _ := cfgUnescape.CompileTemplate("\\\"{{ .version }}\\\"")
	cfgUnescape.Placeholders = []*Placeholder{
		{TemplateText: "\\\"{{ .version }}\\\"", Template: tplJson, Line: 0},
	}
	err = cfgUnescape.ApplyTemplates(dictionary)
	if err != nil {
		t.Fatalf("ApplyTemplates (unescape) failed: %v", err)
	}
	if cfgUnescape.Lines[0] != "\"v1.2.3\"" {
		t.Errorf("expected \"v1.2.3\", got %s", cfgUnescape.Lines[0])
	}
}

// Helper to compile template in tests (simulating what CompilePatterns does)
func (f *ConfigFile) CompileTemplate(tpltxt string) (*template.Template, error) {
	return template.New(f.Name).Funcs(sprig.FuncMap()).Parse(tpltxt)
}

func TestStringMethods(t *testing.T) {
	ph := &Placeholder{
		TemplateText: "{{ .v }}",
		Line:         10,
	}
	if s := ph.String(); s == "" {
		t.Error("Placeholder.String() returned empty string")
	}

	cf := &ConfigFile{
		Name: "test",
		Templates: []*Template{
			{Row: 1, Template: "{{ .v }}"},
		},
		Placeholders: []*Placeholder{ph},
		Lines:        []string{"line1", "line2"},
	}
	// Note: ConfigFile.String() might crash if Placeholders refer to out of bounds Lines
	// but ph.Line is 10 and Lines has 2. Let's fix that for the test.
	ph.Line = 0
	if s := cf.String(); s == "" {
		t.Error("ConfigFile.String() returned empty string")
	}

	cfg := &Config{
		Files: []*ConfigFile{cf},
	}
	if s := cfg.String(); s == "" {
		t.Error("Config.String() returned empty string")
	}
}

func TestCreateConfiguration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "versifyr-create-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	logger := logging.NewLogger()
	
	// Test Simulate
	cfgSim := &Config{
		BasePath: filepath.Join(tmpDir, "sim"),
		Simulate: true,
	}
	if err := cfgSim.CreateConfiguration(logger); err != nil {
		t.Errorf("Simulate CreateConfiguration failed: %v", err)
	}
	if _, err := os.Stat(cfgSim.BasePath); err == nil {
		t.Error("folder created in simulation mode")
	}

	// Test Actual
	cfgAct := &Config{
		BasePath: filepath.Join(tmpDir, "act"),
		Simulate: false,
	}
	if err := cfgAct.CreateConfiguration(logger); err != nil {
		t.Errorf("Actual CreateConfiguration failed: %v", err)
	}
	if _, err := os.Stat(cfgAct.BasePath); os.IsNotExist(err) {
		t.Error("folder not created in actual mode")
	}
	confFile := filepath.Join(cfgAct.BasePath, CONFIG_FILENAME)
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		t.Error("config file not created in actual mode")
	}
}

func TestConfigErrorPaths(t *testing.T) {
	logger := logging.NewLogger()

	// NewConfig non-existent
	os.Setenv(BASEPATH_ENV, "/non/existent/path")
	cfg := &Config{}
	if err := NewConfig(cfg); err == nil {
		t.Error("NewConfig should fail for non-existent path")
	}
	os.Unsetenv(BASEPATH_ENV)

	// CompilePatterns non-existent file
	cfgErr := &Config{
		Files: []*ConfigFile{
			{Path: "/non/existent/file.go"},
		},
	}
	if _, err := cfgErr.CompilePatterns(logger); err == nil {
		t.Error("CompilePatterns should fail for non-existent file")
	}

	// CompilePatterns unclosed template
	tmpDir, _ := os.MkdirTemp("", "versifyr-err-*")
	defer os.RemoveAll(tmpDir)
	fPath := filepath.Join(tmpDir, "unclosed.go")
	os.WriteFile(fPath, []byte("// $versifyr:template={{ .v }}\nconst V = 1"), 0644)
	cfgUnclosed := &Config{
		Files: []*ConfigFile{
			{Path: fPath, Name: "unclosed.go"},
		},
	}
	if _, err := cfgUnclosed.CompilePatterns(logger); err == nil {
		t.Error("CompilePatterns should fail for unclosed template")
	}

	// CompilePatterns invalid template syntax
	fPath2 := filepath.Join(tmpDir, "invalid.go")
	os.WriteFile(fPath2, []byte("// $versifyr:template={{ .v | invalid }}$\nconst V = 1"), 0644)
	cfgInvalid := &Config{
		Files: []*ConfigFile{
			{Path: fPath2, Name: "invalid.go"},
		},
	}
	if _, err := cfgInvalid.CompilePatterns(logger); err == nil {
		t.Error("CompilePatterns should fail for invalid template syntax")
	}
}
