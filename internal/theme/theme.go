package theme

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Theme encapsulates the color palette and styles for UniSocket.
type Theme struct {
	ID          string
	Name        string
	Description string

	// Base palette
	Text        lipgloss.TerminalColor
	Muted       lipgloss.TerminalColor
	Border      lipgloss.TerminalColor
	Highlight   lipgloss.TerminalColor
	TitleText   lipgloss.TerminalColor
	Status      lipgloss.TerminalColor

	// Table styles
	HeaderBg     lipgloss.TerminalColor
	HeaderFg     lipgloss.TerminalColor
	HeaderBorder lipgloss.TerminalColor
	SelectedBg   lipgloss.TerminalColor
	SelectedFg   lipgloss.TerminalColor

	// Modal accents
	ModalBorder lipgloss.TerminalColor
}

// TableStyles generates the table.Styles configuration for the theme.
func (t Theme) TableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(t.HeaderBorder).
		BorderBottom(false).
		Padding(0, 1).
		MarginBottom(1).
		Foreground(t.HeaderFg).
		Background(t.HeaderBg).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(t.SelectedFg).
		Background(t.SelectedBg).
		Bold(true)
	return s
}

// TitleStyle returns the styled banner for the app header.
func (t Theme) TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.TitleText).
		Background(t.Highlight).
		Padding(0, 1)
}

// TableContainerStyle returns the styled border and text container for the table.
func (t Theme) TableContainerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Foreground(t.Text).
		Padding(0, 1)
}

// ThemeBoxStyle returns the styled border container for the theme switcher box.
// It uses t.Border to ensure the border color and style always match the process display box.
func (t Theme) ThemeBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Foreground(t.Text).
		Padding(0, 1)
}

// StatusStyle returns the style for status messages.
func (t Theme) StatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Status)
}

// FooterStyle returns the style for hints and footer controls.
func (t Theme) FooterStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Muted)
}

// ModalStyle returns the styled border container for modals and dialogs.
// It uses t.Border to ensure all dialog and modal borders consistently match the theme's process container.
func (t Theme) ModalStyle() lipgloss.Style {
	border := t.Border
	if border == nil {
		border = t.ModalBorder
	}
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Foreground(t.Text).
		Padding(1, 2)
}

// CommandRowStyle returns the style for interactive input rows (/ and :).
func (t Theme) CommandRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Text).
		Padding(0, 1).
		Bold(true)
}

var allThemes = []Theme{
	{
		ID:          "default",
		Name:        "Default (UniSocket Retro)",
		Description: "Classic cyan & amber retro palette",
		Text:        lipgloss.Color("230"),
		Muted:       lipgloss.Color("245"),
		Border:      lipgloss.Color("31"),
		Highlight:   lipgloss.Color("31"),
		TitleText:   lipgloss.Color("230"),
		Status:      lipgloss.Color("222"),
		HeaderBg:    lipgloss.Color("24"),
		HeaderFg:    lipgloss.Color("229"),
		HeaderBorder: lipgloss.Color("63"),
		SelectedBg:  lipgloss.Color("31"),
		SelectedFg:  lipgloss.Color("230"),
		ModalBorder: lipgloss.Color("31"),
	},
	{
		ID:          "catppuccin",
		Name:        "Catppuccin Mocha",
		Description: "Soothing pastel dark palette",
		Text:        lipgloss.Color("#cdd6f4"),
		Muted:       lipgloss.Color("#6c7086"),
		Border:      lipgloss.Color("#89b4fa"),
		Highlight:   lipgloss.Color("#cba6f7"),
		TitleText:   lipgloss.Color("#11111b"),
		Status:      lipgloss.Color("#fabd2f"),
		HeaderBg:    lipgloss.Color("#313244"),
		HeaderFg:    lipgloss.Color("#cba6f7"),
		HeaderBorder: lipgloss.Color("#89b4fa"),
		SelectedBg:  lipgloss.Color("#cba6f7"),
		SelectedFg:  lipgloss.Color("#11111b"),
		ModalBorder: lipgloss.Color("#89b4fa"),
	},
	{
		ID:          "tokyonight",
		Name:        "Tokyo Night",
		Description: "Downtown Tokyo night neon lights",
		Text:        lipgloss.Color("#c0caf5"),
		Muted:       lipgloss.Color("#565f89"),
		Border:      lipgloss.Color("#7aa2f7"),
		Highlight:   lipgloss.Color("#7aa2f7"),
		TitleText:   lipgloss.Color("#16161e"),
		Status:      lipgloss.Color("#e0af68"),
		HeaderBg:    lipgloss.Color("#24283b"),
		HeaderFg:    lipgloss.Color("#7dcfff"),
		HeaderBorder: lipgloss.Color("#7aa2f7"),
		SelectedBg:  lipgloss.Color("#7aa2f7"),
		SelectedFg:  lipgloss.Color("#16161e"),
		ModalBorder: lipgloss.Color("#7aa2f7"),
	},
	{
		ID:          "gruvbox",
		Name:        "Gruvbox Dark",
		Description: "Warm retro groove with earthy tones",
		Text:        lipgloss.Color("#ebdbb2"),
		Muted:       lipgloss.Color("#928374"),
		Border:      lipgloss.Color("#d65d0e"),
		Highlight:   lipgloss.Color("#fe8019"),
		TitleText:   lipgloss.Color("#282828"),
		Status:      lipgloss.Color("#fabd2f"),
		HeaderBg:    lipgloss.Color("#3c3836"),
		HeaderFg:    lipgloss.Color("#fabd2f"),
		HeaderBorder: lipgloss.Color("#fe8019"),
		SelectedBg:  lipgloss.Color("#fe8019"),
		SelectedFg:  lipgloss.Color("#282828"),
		ModalBorder: lipgloss.Color("#d65d0e"),
	},
	{
		ID:          "dracula",
		Name:        "Dracula",
		Description: "Iconic gothic dark theme with vibrant accents",
		Text:        lipgloss.Color("#f8f8f2"),
		Muted:       lipgloss.Color("#6272a4"),
		Border:      lipgloss.Color("#bd93f9"),
		Highlight:   lipgloss.Color("#bd93f9"),
		TitleText:   lipgloss.Color("#282a36"),
		Status:      lipgloss.Color("#ffb86c"),
		HeaderBg:    lipgloss.Color("#44475a"),
		HeaderFg:    lipgloss.Color("#50fa7b"),
		HeaderBorder: lipgloss.Color("#bd93f9"),
		SelectedBg:  lipgloss.Color("#bd93f9"),
		SelectedFg:  lipgloss.Color("#282a36"),
		ModalBorder: lipgloss.Color("#bd93f9"),
	},
	{
		ID:          "nord",
		Name:        "Nord",
		Description: "Arctic, north-bluish clean palette",
		Text:        lipgloss.Color("#eceff4"),
		Muted:       lipgloss.Color("#4c566a"),
		Border:      lipgloss.Color("#88c0d0"),
		Highlight:   lipgloss.Color("#88c0d0"),
		TitleText:   lipgloss.Color("#2e3440"),
		Status:      lipgloss.Color("#ebcb8b"),
		HeaderBg:    lipgloss.Color("#3b4252"),
		HeaderFg:    lipgloss.Color("#8fbcbb"),
		HeaderBorder: lipgloss.Color("#81a1c1"),
		SelectedBg:  lipgloss.Color("#88c0d0"),
		SelectedFg:  lipgloss.Color("#2e3440"),
		ModalBorder: lipgloss.Color("#88c0d0"),
	},
	{
		ID:          "onedark",
		Name:        "One Dark",
		Description: "Atom & Neovim balanced classic",
		Text:        lipgloss.Color("#abb2bf"),
		Muted:       lipgloss.Color("#5c6370"),
		Border:      lipgloss.Color("#61afef"),
		Highlight:   lipgloss.Color("#61afef"),
		TitleText:   lipgloss.Color("#282c34"),
		Status:      lipgloss.Color("#e5c07b"),
		HeaderBg:    lipgloss.Color("#353b45"),
		HeaderFg:    lipgloss.Color("#98c379"),
		HeaderBorder: lipgloss.Color("#61afef"),
		SelectedBg:  lipgloss.Color("#61afef"),
		SelectedFg:  lipgloss.Color("#282c34"),
		ModalBorder: lipgloss.Color("#61afef"),
	},
	{
		ID:          "rose-pine",
		Name:        "Rosé Pine",
		Description: "Soho vibes, pine tones, and cozy aesthetic",
		Text:        lipgloss.Color("#e0def4"),
		Muted:       lipgloss.Color("#6e6a86"),
		Border:      lipgloss.Color("#ebbcba"),
		Highlight:   lipgloss.Color("#ebbcba"),
		TitleText:   lipgloss.Color("#191724"),
		Status:      lipgloss.Color("#f6c177"),
		HeaderBg:    lipgloss.Color("#26233a"),
		HeaderFg:    lipgloss.Color("#9ccfd8"),
		HeaderBorder: lipgloss.Color("#ebbcba"),
		SelectedBg:  lipgloss.Color("#ebbcba"),
		SelectedFg:  lipgloss.Color("#191724"),
		ModalBorder: lipgloss.Color("#ebbcba"),
	},
	{
		ID:          "solarized-dark",
		Name:        "Solarized Dark",
		Description: "Engineered low-contrast terminal legend",
		Text:        lipgloss.Color("#93a1a1"),
		Muted:       lipgloss.Color("#586e75"),
		Border:      lipgloss.Color("#268bd2"),
		Highlight:   lipgloss.Color("#268bd2"),
		TitleText:   lipgloss.Color("#002b36"),
		Status:      lipgloss.Color("#b58900"),
		HeaderBg:    lipgloss.Color("#073642"),
		HeaderFg:    lipgloss.Color("#2aa198"),
		HeaderBorder: lipgloss.Color("#268bd2"),
		SelectedBg:  lipgloss.Color("#268bd2"),
		SelectedFg:  lipgloss.Color("#002b36"),
		ModalBorder: lipgloss.Color("#268bd2"),
	},
	{
		ID:          "monokai",
		Name:        "Monokai Pro",
		Description: "High-contrast vibrant developer theme",
		Text:        lipgloss.Color("#fcfcfa"),
		Muted:       lipgloss.Color("#727072"),
		Border:      lipgloss.Color("#ff6188"),
		Highlight:   lipgloss.Color("#ffd866"),
		TitleText:   lipgloss.Color("#2d2a2e"),
		Status:      lipgloss.Color("#fc9867"),
		HeaderBg:    lipgloss.Color("#403e41"),
		HeaderFg:    lipgloss.Color("#a9dc76"),
		HeaderBorder: lipgloss.Color("#ffd866"),
		SelectedBg:  lipgloss.Color("#ffd866"),
		SelectedFg:  lipgloss.Color("#2d2a2e"),
		ModalBorder: lipgloss.Color("#ff6188"),
	},
	{
		ID:          "cyberpunk",
		Name:        "Cyberpunk Neon",
		Description: "Neon synthwave glow for night owls",
		Text:        lipgloss.Color("#ffffff"),
		Muted:       lipgloss.Color("#846b8a"),
		Border:      lipgloss.Color("#00ffff"),
		Highlight:   lipgloss.Color("#ff007f"),
		TitleText:   lipgloss.Color("#ffffff"),
		Status:      lipgloss.Color("#ffe600"),
		HeaderBg:    lipgloss.Color("#3f1052"),
		HeaderFg:    lipgloss.Color("#ffe600"),
		HeaderBorder: lipgloss.Color("#00ffff"),
		SelectedBg:  lipgloss.Color("#ff007f"),
		SelectedFg:  lipgloss.Color("#ffffff"),
		ModalBorder: lipgloss.Color("#00ffff"),
	},
}

// List returns a slice of all available themes.
func List() []Theme {
	themes := make([]Theme, len(allThemes))
	copy(themes, allThemes)
	return themes
}

// Default returns the default UniSocket theme.
func Default() Theme {
	return allThemes[0]
}

// Register registers a custom or external theme, replacing any existing theme with the same ID.
func Register(th Theme) {
	if th.ID == "" {
		return
	}
	for i, existing := range allThemes {
		if existing.ID == th.ID {
			allThemes[i] = th
			return
		}
	}
	allThemes = append(allThemes, th)
}

// Get resolves a theme by ID, name, or common aliases.
func Get(query string) (Theme, bool) {
	clean := strings.ToLower(strings.TrimSpace(query))
	clean = strings.ReplaceAll(clean, "_", "-")
	clean = strings.ReplaceAll(clean, " ", "-")

	// Direct ID match
	for _, th := range allThemes {
		if th.ID == clean || strings.ToLower(th.Name) == clean {
			return th, true
		}
	}

	// Alias matching
	switch clean {
	case "cattputchin", "catppuchin", "catppuccin-mocha", "mocha":
		return Get("catppuccin")
	case "tokyo", "tokyo-night":
		return Get("tokyonight")
	case "gruv", "gruvbox-dark":
		return Get("gruvbox")
	case "one-dark", "atom":
		return Get("onedark")
	case "rose", "rosepine":
		return Get("rose-pine")
	case "solarized", "solarizedark":
		return Get("solarized-dark")
	case "monokaipro", "monokai-pro":
		return Get("monokai")
	case "synthwave", "neon":
		return Get("cyberpunk")
	case "retro":
		return Get("default")
	}

	// Partial prefix / substring match
	for _, th := range allThemes {
		if strings.HasPrefix(th.ID, clean) || strings.HasPrefix(strings.ToLower(th.Name), clean) {
			return th, true
		}
	}

	return Theme{}, false
}
