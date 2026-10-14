package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Senasphy/unisocket/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// Config represents the application configuration.
type Config struct {
	filePath string `json:"-"`

	// Active color theme (e.g., "default", "catppuccin", "tokyonight", "gruvbox", etc.)
	Theme string `json:"theme"`

	// Refresh interval in seconds for live connection scanning (default: 2)
	RefreshInterval int `json:"refresh_interval"`

	// Default sort mode on launch: "" (none), "name", "process-name", "pid"
	DefaultSort string `json:"default_sort"`

	// Default filename template for snapshot exports (default: "portList.json")
	SnapshotFilename string `json:"snapshot_filename"`

	// Default startup connection filters
	Filters FiltersConfig `json:"filters"`

	// Custom user-defined themes
	CustomThemes []CustomTheme `json:"custom_themes,omitempty"`
}

// FiltersConfig defines the startup filter toggles.
type FiltersConfig struct {
	TCP         bool   `json:"tcp"`
	UDP         bool   `json:"udp"`
	Established bool   `json:"established"`
	Listening   bool   `json:"listening"`
	FindByName  string `json:"find_by_name"`
}

// CustomTheme allows defining new themes or tweaking colors via config.
type CustomTheme struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Text      string `json:"text"`
	Muted     string `json:"muted"`
	Border    string `json:"border"`
	Highlight string `json:"highlight"`
	TitleText string `json:"title_text"`
	Status    string `json:"status"`

	HeaderBg     string `json:"header_bg"`
	HeaderFg     string `json:"header_fg"`
	HeaderBorder string `json:"header_border"`
	SelectedBg   string `json:"selected_bg"`
	SelectedFg   string `json:"selected_fg"`
	ModalBorder  string `json:"modal_border"`
}

// ToTheme converts a CustomTheme config into a theme.Theme.
func (ct CustomTheme) ToTheme() theme.Theme {
	colorOr := func(val, def string) lipgloss.TerminalColor {
		val = strings.TrimSpace(val)
		if val == "" {
			return lipgloss.Color(def)
		}
		return lipgloss.Color(val)
	}

	name := ct.Name
	if name == "" {
		name = ct.ID
	}

	borderDefault := ct.Border
	if borderDefault == "" {
		borderDefault = "#89b4fa"
	}

	return theme.Theme{
		ID:           ct.ID,
		Name:         name,
		Description:  ct.Description,
		Text:         colorOr(ct.Text, "#cdd6f4"),
		Muted:        colorOr(ct.Muted, "#6c7086"),
		Border:       colorOr(ct.Border, "#89b4fa"),
		Highlight:    colorOr(ct.Highlight, "#cba6f7"),
		TitleText:    colorOr(ct.TitleText, "#11111b"),
		Status:       colorOr(ct.Status, "#fab387"),
		HeaderBg:     colorOr(ct.HeaderBg, "#313244"),
		HeaderFg:     colorOr(ct.HeaderFg, "#cba6f7"),
		HeaderBorder: colorOr(ct.HeaderBorder, "#89b4fa"),
		SelectedBg:   colorOr(ct.SelectedBg, "#cba6f7"),
		SelectedFg:   colorOr(ct.SelectedFg, "#11111b"),
		ModalBorder:  colorOr(ct.ModalBorder, borderDefault),
	}
}

// DefaultPath returns the default configuration file location.
func DefaultPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", "unisocket", "config.json")
	}
	return filepath.Join(configDir, "unisocket", "config.json")
}

// DefaultConfig returns a new Config with standard default values.
func DefaultConfig() *Config {
	return &Config{
		filePath:         DefaultPath(),
		Theme:            "default",
		RefreshInterval:  2,
		DefaultSort:      "",
		SnapshotFilename: "portList.json",
		Filters: FiltersConfig{
			TCP:         false,
			UDP:         false,
			Established: false,
			Listening:   false,
			FindByName:  "",
		},
		CustomThemes: []CustomTheme{},
	}
}

// Load reads the configuration from customPath or the default path.
// If the configuration file does not exist, it creates one with default values.
func Load(customPath string) *Config {
	targetPath := strings.TrimSpace(customPath)
	if targetPath == "" {
		targetPath = DefaultPath()
	}

	cfg := DefaultConfig()
	cfg.filePath = targetPath

	data, err := os.ReadFile(targetPath)
	if err != nil {
		// Auto-generate default configuration file if not found
		_ = cfg.Save()
		return cfg
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg
	}

	cfg.filePath = targetPath

	// Validate bounds
	if cfg.RefreshInterval <= 0 {
		cfg.RefreshInterval = 2
	}
	if strings.TrimSpace(cfg.SnapshotFilename) == "" {
		cfg.SnapshotFilename = "portList.json"
	}
	if strings.TrimSpace(cfg.Theme) == "" {
		cfg.Theme = "default"
	}

	// Register any custom themes defined in the config file
	for _, ct := range cfg.CustomThemes {
		if ct.ID != "" {
			theme.Register(ct.ToTheme())
		}
	}

	return cfg
}

// FilePath returns the location of the config file.
func (c *Config) FilePath() string {
	if c.filePath == "" {
		return DefaultPath()
	}
	return c.filePath
}

// Save writes the current configuration to disk as formatted JSON.
func (c *Config) Save() error {
	path := c.FilePath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

// SaveTheme updates the theme in the default config file and persists it.
func SaveTheme(themeID string) error {
	cfg := Load("")
	cfg.Theme = themeID
	return cfg.Save()
}
