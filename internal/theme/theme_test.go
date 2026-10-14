package theme

import (
	"testing"
)

func TestThemesList(t *testing.T) {
	themes := List()
	if len(themes) < 10 {
		t.Fatalf("expected at least 10 themes, got %d", len(themes))
	}

	def := Default()
	if def.ID != "default" {
		t.Fatalf("expected default theme ID 'default', got %s", def.ID)
	}
}

func TestThemeGet(t *testing.T) {
	tests := []struct {
		query    string
		expected string
		found    bool
	}{
		{"default", "default", true},
		{"catppuccin", "catppuccin", true},
		{"cattputchin", "catppuccin", true},
		{"tokyonight", "tokyonight", true},
		{"tokyo", "tokyonight", true},
		{"gruv", "gruvbox", true},
		{"gruvbox", "gruvbox", true},
		{"dracula", "dracula", true},
		{"nord", "nord", true},
		{"onedark", "onedark", true},
		{"rose-pine", "rose-pine", true},
		{"solarized-dark", "solarized-dark", true},
		{"solarized", "solarized-dark", true},
		{"monokai", "monokai", true},
		{"cyberpunk", "cyberpunk", true},
		{"nonexistent-theme", "", false},
	}

	for _, tc := range tests {
		th, ok := Get(tc.query)
		if ok != tc.found {
			t.Errorf("Get(%q) ok = %v, expected %v", tc.query, ok, tc.found)
		}
		if ok && th.ID != tc.expected {
			t.Errorf("Get(%q) theme.ID = %q, expected %q", tc.query, th.ID, tc.expected)
		}
	}
}

func TestThemeStyles(t *testing.T) {
	for _, th := range List() {
		styles := th.TableStyles()
		if styles.Header.GetBold() != true {
			t.Errorf("theme %s table header should be bold", th.ID)
		}
		if styles.Selected.GetBold() != true {
			t.Errorf("theme %s selected row should be bold", th.ID)
		}
		_ = th.TitleStyle().Render("test")
		_ = th.TableContainerStyle().Render("test")
		_ = th.StatusStyle().Render("test")
		_ = th.FooterStyle().Render("test")
		_ = th.ModalStyle().Render("test")
		_ = th.CommandRowStyle().Render("test")
	}
}

func TestRegisterTheme(t *testing.T) {
	custom := Theme{
		ID:          "my-custom",
		Name:        "My Custom Theme",
		Description: "A custom test theme",
	}
	Register(custom)

	found, ok := Get("my-custom")
	if !ok {
		t.Fatalf("expected registered custom theme to be found")
	}
	if found.Name != "My Custom Theme" {
		t.Fatalf("expected custom theme name 'My Custom Theme', got %s", found.Name)
	}
}
