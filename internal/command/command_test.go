package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
	"github.com/zukrin/versifyr/internal/configuration"
	"github.com/zukrin/versifyr/internal/logging"
)

func TestCLIWorkflow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "versifyr-integration-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Setup App (similar to main.go but for testing)
	cfg := &configuration.Config{}
	logger := logging.NewLogger()
	// Disable pterm output for tests if possible or just let it write to default writers
	
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
				configuration.NewConfig(cfg)
				cfg.CompilePatterns(logger)
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
	os.WriteFile(testFile, []byte(testContent), 0644)

	// Update configuration to include this file
	// Actually, the default configuration created by init includes version.go in internal/versifyr/version.go
	// Let's overwrite configuration.yaml to match our test file
	customConfig := `
files:
  - name: version.go
    type: go
    path: version.go
`
	os.WriteFile(confFile, []byte(customConfig), 0644)

	// 3. Run Show (mostly to see if it doesn't crash)
	err = app.Run([]string{"versifyr-test", "show"})
	if err != nil {
		t.Fatalf("show failed: %v", err)
	}

	// 4. Run Set
	err = app.Run([]string{"versifyr-test", "set", "version=v1.2.3"})
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	// Verify change
	updatedContent, _ := os.ReadFile(testFile)
	if !strings.Contains(string(updatedContent), "const Version = \"v1.2.3\"") {
		t.Errorf("file was not updated correctly. Got:\n%s", string(updatedContent))
	}
}

func TestCLIWorkflowErrors(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "versifyr-errors-*")
	defer os.RemoveAll(tmpDir)
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	app := &cli.App{
		Metadata: map[string]interface{}{
			"config": &configuration.Config{BasePath: ".versifyr"},
			"logger": logging.NewLogger(),
		},
		Commands: []*cli.Command{InitCommand, SetCommand},
	}

	// 1. Init when folder exists
	os.Mkdir(".versifyr", 0755)
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
}
