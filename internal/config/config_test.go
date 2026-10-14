package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Senasphy/unisocket/internal/theme"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Theme != "default" {
		t.Fatalf("expected default theme, got %s", cfg.Theme)
	}
	if cfg.RefreshInterval != 2 {
		t.Fatalf("expected refresh interval 2, got %d", cfg.RefreshInterval)
	}
	if cfg.SnapshotFilename != "portList.json" {
		t.Fatalf("expected snapshot filename portList.json, got %s", cfg.SnapshotFilename)
	}
}

func TestLoadAutoCreate(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "unisocket", "config.json")

	cfg := Load(configPath)
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	// Verify file was automatically created
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected config file to be auto-created at %s: %v", configPath, err)
	}

	if cfg.Theme != "default" {
		t.Fatalf("expected theme 'default', got %s", cfg.Theme)
	}
}

func TestSaveAndReloadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "custom-config.json")

	cfg := Load(configPath)
	cfg.Theme = "tokyonight"
	cfg.RefreshInterval = 5
	cfg.DefaultSort = "pid"
	cfg.Filters.TCP = true
	cfg.Filters.Listening = true

	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	reloaded := Load(configPath)
	if reloaded.Theme != "tokyonight" {
		t.Fatalf("expected theme 'tokyonight', got %s", reloaded.Theme)
	}
	if reloaded.RefreshInterval != 5 {
		t.Fatalf("expected refresh interval 5, got %d", reloaded.RefreshInterval)
	}
	if reloaded.DefaultSort != "pid" {
		t.Fatalf("expected default sort 'pid', got %s", reloaded.DefaultSort)
	}
	if !reloaded.Filters.TCP || !reloaded.Filters.Listening {
		t.Fatalf("expected filters to be preserved")
	}
}

func TestCustomThemeRegistration(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "custom-themes.json")

	jsonContent := `{
		"theme": "matrix",
		"custom_themes": [
			{
				"id": "matrix",
				"name": "The Matrix",
				"description": "Green terminal phosphor",
				"text": "#00ff00",
				"highlight": "#003300",
				"border": "#00ff00"
			}
		]
	}`

	if err := os.WriteFile(configPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("failed to write custom theme config: %v", err)
	}

	cfg := Load(configPath)
	if cfg.Theme != "matrix" {
		t.Fatalf("expected theme 'matrix', got %s", cfg.Theme)
	}

	matrixTheme, ok := theme.Get("matrix")
	if !ok {
		t.Fatalf("expected custom theme 'matrix' to be registered in theme package")
	}
	if matrixTheme.Name != "The Matrix" {
		t.Fatalf("expected theme name 'The Matrix', got %s", matrixTheme.Name)
	}
}
