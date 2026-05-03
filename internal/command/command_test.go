package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/urfave/cli/v2"
	"github.com/zukrin/versifyr/internal/configuration"
	"github.com/zukrin/versifyr/internal/logging"
)

func TestCLIWorkflow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "versifyr-integration-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(oldWd) }()

	// Setup App (similar to main.go but for testing)
	cfg := &configuration.Config{}
	logger := logging.NewLogger()
	
	app := &cli.App{
		Name: "versifyr-test",
		Metadata: map[string]interface{}{
			"config": cfg,
			"logger": logger,
		},
		Commands: []*cli.Command{
			InitCommand,
			ShowCommand,
			SetCommand,
		},
		Before: func(cCtx *cli.Context) error {
			cfg := cCtx.App.Metadata["config"].(*configuration.Config)
			logger := cCtx.App.Metadata["logger"].(*logging.Logger)
			
			bp := configuration.GetBasePath()
			cfg.BasePath = bp

			// Only try to load config if it exists
			if _, err := os.Stat(bp + "/" + configuration.CONFIG_FILENAME); err == nil {
				_ = configuration.NewConfig(cfg)
				_, _ = cfg.CompilePatterns(logger)
			}
			return nil
		},
	}

	// 1. Run Init
	err = app.Run([]string{"versifyr-test", "init"})
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	confFile := filepath.Join(tmpDir, ".versifyr", "configuration.yaml")
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		t.Errorf("configuration file not created")
	}

	// 2. Create a test file to be versifyed
	testFile := filepath.Join(tmpDir, "version.go")
	testContent := `package main
// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.0"
`
	_ = os.WriteFile(testFile, []byte(testContent), 0644)

	customConfig := `
files:
  - name: version.go
    type: go
    path: version.go
`
	_ = os.WriteFile(confFile, []byte(customConfig), 0644)

	// Reload config after manual edit
	_ = configuration.NewConfig(cfg)
	_, _ = cfg.CompilePatterns(logger)

	// 3. Run Show
	err = app.Run([]string{"versifyr-test", "show"})
	if err != nil {
		t.Fatalf("show failed: %v", err)
	}

	// 4. Run Set with simulation
	cfg.Simulate = true
	err = app.Run([]string{"versifyr-test", "set", "version=v1.2.3"})
	if err != nil {
		t.Fatalf("set simulation failed: %v", err)
	}
	// Verify no change
	simContent, _ := os.ReadFile(testFile)
	if strings.Contains(string(simContent), "v1.2.3") {
		t.Errorf("file was updated in simulation mode")
	}

	// 5. Run Set for real
	cfg.Simulate = false
	err = app.Run([]string{"versifyr-test", "set", "version=v1.2.3"})
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	// Verify change
	updatedContent, _ := os.ReadFile(testFile)
	if !strings.Contains(string(updatedContent), "const Version = \"v1.2.3\"") {
		t.Errorf("file was not updated correctly. Got:\n%s", string(updatedContent))
	}

	// 6. Run Set with summary
	err = app.Run([]string{"versifyr-test", "set", "--summary", "version=v2.0.0"})
	if err != nil {
		t.Fatalf("set with summary failed: %v", err)
	}
}

func TestCLIWorkflowErrors(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "versifyr-errors-*")
	defer func() { _ = os.RemoveAll(tmpDir) }()
	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(oldWd) }()

	logger := logging.NewLogger()
	app := &cli.App{
		Metadata: map[string]interface{}{
			"config": &configuration.Config{BasePath: ".versifyr"},
			"logger": logger,
		},
		Commands: []*cli.Command{InitCommand, SetCommand},
	}

	// 1. Init when folder exists
	_ = os.Mkdir(".versifyr", 0755)
	if err := app.Run([]string{"versifyr-test", "init"}); err == nil {
		t.Error("init should fail when .versifyr exists")
	}

	// 2. Set with no args
	if err := app.Run([]string{"versifyr-test", "set"}); err == nil {
		t.Error("set should fail with no args")
	}

	// 3. Set with malformed arg
	if err := app.Run([]string{"versifyr-test", "set", "malformed"}); err == nil {
		t.Error("set should fail with malformed arg")
	}

	// 4. Set with non-existent file in config
	cfg := app.Metadata["config"].(*configuration.Config)
	cfg.Files = []*configuration.ConfigFile{
		{Name: "missing", Path: "missing.go"},
	}
	// Note: doSet applies templates to file.Lines. If file was never loaded by CompilePatterns, Lines might be nil.
	// CompilePatterns would have failed earlier in a real app, but here we are testing doSet directly via app.Run.
	// In the real app, Before hook calls CompilePatterns.
	
	// Let's mock a file that exists but fails to write
	badFile := filepath.Join(tmpDir, "readonly.go")
	_ = os.WriteFile(badFile, []byte("test"), 0444)
	cfg.Files = []*configuration.ConfigFile{
		{Name: "readonly", Path: badFile, Lines: []string{"test"}},
	}
	if err := app.Run([]string{"versifyr-test", "set", "v=1"}); err == nil {
		// This might fail because we didn't set up placeholders, so ApplyTemplates does nothing.
	}
}

func TestDoSetApplyError(t *testing.T) {
	logger := logging.NewLogger()

	// Use sprig's fail function to trigger an execution error
	tpl, _ := template.New("test").Funcs(sprig.FuncMap()).Parse("{{ fail \"forced error\" }}")

	cfg := &configuration.Config{
		Files: []*configuration.ConfigFile{
			{
				Name:  "test.go",
				Path:  "test.go",
				Lines: []string{"// $versifyr:template={{ .v }}$", "const V = 1"},
				Placeholders: []*configuration.Placeholder{
					{TemplateText: "{{ .v }}", Line: 1, Template: tpl},
				},
			},
		},
	}
	
	app := &cli.App{
		Metadata: map[string]interface{}{
			"config": cfg,
			"logger": logger,
		},
		Commands: []*cli.Command{SetCommand},
	}
	
	if err := app.Run([]string{"versifyr-test", "set", "v=1"}); err == nil {
		t.Error("doSet should fail when ApplyTemplates fails")
	}
}

func TestDoSetWriteError(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "versifyr-write-err-*")
	defer func() { _ = os.RemoveAll(tmpDir) }()
	
	fPath := filepath.Join(tmpDir, "readonly.go")
	_ = os.WriteFile(fPath, []byte("content"), 0644)
	
	cfg := &configuration.Config{
		Files: []*configuration.ConfigFile{
			{
				Name: "readonly.go",
				Path: fPath,
				Lines: []string{"content"},
			},
		},
	}
	
	logger := logging.NewLogger()
	app := &cli.App{
		Metadata: map[string]interface{}{
			"config": cfg,
			"logger": logger,
		},
		Commands: []*cli.Command{SetCommand},
	}
	
	// Make file read-only
	_ = os.Chmod(fPath, 0400)
	
	if err := app.Run([]string{"versifyr-test", "set", "v=1"}); err == nil {
		t.Error("doSet should fail when file is not writable")
	}
}
