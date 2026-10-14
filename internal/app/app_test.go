package app

import (
	"strings"
	"testing"

	"github.com/Senasphy/unisocket/internal/config"
	"github.com/Senasphy/unisocket/internal/connections"
	"github.com/Senasphy/unisocket/internal/theme"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newTestModel(th ...theme.Theme) Model {
	columns := []table.Column{
		{Title: "PORT", Width: 10},
		{Title: "PROCESS ID", Width: 10},
		{Title: "PROCESS NAME", Width: 10},
		{Title: "PROTOCOL", Width: 10},
		{Title: "STATE", Width: 10},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{
			{"8080", "1234", "node", "tcp", "LISTEN"},
			{"3000", "5678", "go", "tcp", "ESTABLISHED"},
		}),
	)
	m := New(t, connections.Filters{}, nil, th...)
	return m.(Model)
}

func TestModelThemeInit(t *testing.T) {
	tokyo, _ := theme.Get("tokyonight")
	m := newTestModel(tokyo)

	if m.currentTheme.ID != "tokyonight" {
		t.Fatalf("expected theme tokyonight, got %s", m.currentTheme.ID)
	}

	view := m.View()
	if !strings.Contains(strings.ToLower(view), "tokyonight") {
		t.Fatalf("expected view to display active theme, got: %s", view)
	}
}

func TestThemePickerNavigationAndApply(t *testing.T) {
	m := newTestModel()

	// Press 't' to open theme picker
	mModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m = mModel.(Model)

	if !m.isThemePicker {
		t.Fatalf("expected isThemePicker to be true after pressing 't'")
	}

	// Press 'down' or 'j' to navigate and preview
	mModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = mModel.(Model)

	if m.themePickerIndex != 1 {
		t.Fatalf("expected themePickerIndex 1, got %d", m.themePickerIndex)
	}
	expectedTheme := theme.List()[1]
	if m.currentTheme.ID != expectedTheme.ID {
		t.Fatalf("expected live preview theme %s, got %s", expectedTheme.ID, m.currentTheme.ID)
	}

	// Press Enter to confirm
	mModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mModel.(Model)

	if m.isThemePicker {
		t.Fatalf("expected isThemePicker to be false after pressing Enter")
	}
	if m.currentTheme.ID != expectedTheme.ID {
		t.Fatalf("expected applied theme %s, got %s", expectedTheme.ID, m.currentTheme.ID)
	}
}

func TestThemePickerCancel(t *testing.T) {
	def := theme.Default()
	m := newTestModel(def)

	// Open picker
	mModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m = mModel.(Model)

	// Move cursor
	mModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = mModel.(Model)

	// Cancel with Esc
	mModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = mModel.(Model)

	if m.isThemePicker {
		t.Fatalf("expected isThemePicker to be false after pressing Esc")
	}
	if m.currentTheme.ID != def.ID {
		t.Fatalf("expected reverted theme %s, got %s", def.ID, m.currentTheme.ID)
	}
}

func TestThemeCommand(t *testing.T) {
	m := newTestModel()

	// Run :theme catppuccin
	m.applyCommand("theme catppuccin")
	if m.currentTheme.ID != "catppuccin" {
		t.Fatalf("expected theme catppuccin, got %s", m.currentTheme.ID)
	}

	// Run alias :t gruv
	m.applyCommand("t gruv")
	if m.currentTheme.ID != "gruvbox" {
		t.Fatalf("expected theme gruvbox, got %s", m.currentTheme.ID)
	}

	// Run :theme to open picker
	m.applyCommand("theme")
	if !m.isThemePicker {
		t.Fatalf("expected :theme to open theme picker modal")
	}
}

func TestThemePickerSideBySide(t *testing.T) {
	m := newTestModel()
	m.width = 160 // wide screen
	m.openThemePicker()

	view := m.View()
	if !strings.Contains(view, "Themes (Preview)") {
		t.Fatalf("expected side-by-side theme title 'Themes (Preview)' in view, got: %s", view)
	}
}

func TestThemePickerWrapped(t *testing.T) {
	m := newTestModel()
	m.width = 70 // smaller screen
	m.openThemePicker()

	view := m.View()
	if !strings.Contains(view, "Theme Switcher (Live Preview)") {
		t.Fatalf("expected wrapped theme title 'Theme Switcher (Live Preview)' in view, got: %s", view)
	}
}

func TestConfigOptionsInApp(t *testing.T) {
	columns := []table.Column{
		{Title: "PORT", Width: 10},
		{Title: "PROCESS ID", Width: 10},
		{Title: "PROCESS NAME", Width: 10},
		{Title: "PROTOCOL", Width: 10},
		{Title: "STATE", Width: 10},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{
			{"8080", "1234", "node", "tcp", "LISTEN"},
			{"3000", "5678", "go", "tcp", "ESTABLISHED"},
		}),
	)

	cfg := &config.Config{
		Theme:            "nord",
		RefreshInterval:  4,
		DefaultSort:      "pid",
		SnapshotFilename: "custom-ports.json",
	}

	mModel := New(tbl, connections.Filters{}, cfg)
	m := mModel.(Model)

	if m.currentTheme.ID != "nord" {
		t.Fatalf("expected theme 'nord', got %s", m.currentTheme.ID)
	}
	if m.sortMode != "pid" {
		t.Fatalf("expected sortMode 'pid', got %s", m.sortMode)
	}
}

func TestThemeBoxBorderConsistency(t *testing.T) {
	for _, th := range theme.List() {
		boxBorder := th.ThemeBoxStyle().GetBorderTopForeground()
		tableBorder := th.TableContainerStyle().GetBorderTopForeground()
		if boxBorder != tableBorder {
			t.Errorf("theme %s: box border color %v does not match table border color %v", th.ID, boxBorder, tableBorder)
		}
	}
}

func TestSearchExitWithoutQuitting(t *testing.T) {
	m := newTestModel()

	// 1. Enter search mode
	mModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = mModel.(Model)
	if !m.isSearching {
		t.Fatalf("expected isSearching to be true")
	}

	// 2. Press Esc while typing in search - should exit search mode without quitting
	mModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = mModel.(Model)
	if m.isSearching {
		t.Fatalf("expected isSearching to be false after Esc")
	}
	if cmd != nil {
		t.Fatalf("expected nil cmd after Esc in search mode, got non-nil (quitting)")
	}

	// 3. Re-enter search and test Ctrl+C - should exit search without quitting
	mModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = mModel.(Model)
	mModel, cmd = m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	m = mModel.(Model)
	if m.isSearching {
		t.Fatalf("expected isSearching to be false after Ctrl+C")
	}
	if cmd != nil {
		t.Fatalf("expected nil cmd after Ctrl+C in search mode, got non-nil")
	}

	// 4. Set search query via Enter and then test Esc clears search without quitting
	mModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = mModel.(Model)
	m.searchInput.SetValue("node")
	m.searchQuery = "node"
	mModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mModel.(Model)
	if m.searchQuery != "node" {
		t.Fatalf("expected search query 'node'")
	}

	// Press Esc while on main table with active filter -> should clear query, not quit
	mModel, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = mModel.(Model)
	if m.searchQuery != "" {
		t.Fatalf("expected search query to be cleared after Esc")
	}
	if cmd != nil {
		t.Fatalf("expected Esc to clear search without quitting")
	}
}

func TestBottomDeckExactWidth(t *testing.T) {
	// Test side-by-side mode (theme picker open on wide terminal)
	m := newTestModel()
	m.width = 160
	m.height = 40
	m.updateTableDimensions()
	m.openThemePicker()

	tableView := m.currentTheme.TableContainerStyle().Render(m.table.View())
	tableWidth := lipgloss.Width(tableView)
	themeBox := m.renderThemePickerBox(tableView, true, 0)
	themeWidth := lipgloss.Width(themeBox)
	totalWidth := tableWidth + gapWidth + themeWidth

	deck := m.renderBottomDeck(totalWidth, tableWidth, true)
	if lipgloss.Width(deck) != totalWidth {
		t.Fatalf("Side-by-side: deck width (%d) != totalWidth (%d)", lipgloss.Width(deck), totalWidth)
	}

	searchBox := m.renderSearchBox(totalWidth)
	if lipgloss.Width(searchBox) != totalWidth {
		t.Fatalf("Search box width (%d) != totalWidth (%d)", lipgloss.Width(searchBox), totalWidth)
	}

	commandBox := m.renderCommandBox(totalWidth)
	if lipgloss.Width(commandBox) != totalWidth {
		t.Fatalf("Command box width (%d) != totalWidth (%d)", lipgloss.Width(commandBox), totalWidth)
	}

	// Test normal mode (no theme picker)
	m2 := newTestModel()
	m2.width = 120
	m2.height = 30
	m2.updateTableDimensions()

	tv2 := m2.currentTheme.TableContainerStyle().Render(m2.table.View())
	tw2 := lipgloss.Width(tv2)

	deck2 := m2.renderBottomDeck(tw2, tw2, false)
	if lipgloss.Width(deck2) != tw2 {
		t.Fatalf("Normal mode: deck width (%d) != tableWidth (%d)", lipgloss.Width(deck2), tw2)
	}
}

func TestResponsiveResize(t *testing.T) {
	m := newTestModel()

	// Resize to 80x24 (standard small terminal)
	mModel, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = mModel.(Model)

	if m.width != 80 || m.height != 24 {
		t.Fatalf("expected 80x24, got %dx%d", m.width, m.height)
	}
	if m.isTooSmall() {
		t.Fatalf("80x24 should not be marked as too small")
	}

	tv80 := m.currentTheme.TableContainerStyle().Render(m.table.View())
	w80 := lipgloss.Width(tv80)
	if w80 > 80-pagePaddingH {
		t.Fatalf("table container width %d exceeds available terminal space %d", w80, 80-pagePaddingH)
	}

	// Resize to 140x45 (wide terminal)
	mModel, _ = m.Update(tea.WindowSizeMsg{Width: 140, Height: 45})
	m = mModel.(Model)

	tv140 := m.currentTheme.TableContainerStyle().Render(m.table.View())
	w140 := lipgloss.Width(tv140)
	if w140 <= w80 {
		t.Fatalf("expected table width at 140 (%d) to expand wider than at 80 (%d)", w140, w80)
	}
	if w140 > 140-pagePaddingH {
		t.Fatalf("table container width %d exceeds available terminal space %d", w140, 140-pagePaddingH)
	}
}

func TestScreenSizeGuard(t *testing.T) {
	m := newTestModel()

	// Terminal too narrow (< 64)
	mModel, _ := m.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
	m = mModel.(Model)
	if !m.isTooSmall() {
		t.Fatalf("expected 50x20 to be marked as too small")
	}
	view := m.View()
	if !strings.Contains(view, "Terminal Window Too Small") {
		t.Fatalf("expected view to display too small message, got: %s", view)
	}

	// Quitting from too small screen
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatalf("expected 'q' to quit when screen is too small")
	}

	// Terminal too short (< 12)
	mModel, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 10})
	m = mModel.(Model)
	if !m.isTooSmall() {
		t.Fatalf("expected 100x10 to be marked as too small")
	}

	// Enlarge terminal back to sufficient size
	mModel, _ = m.Update(tea.WindowSizeMsg{Width: 90, Height: 25})
	m = mModel.(Model)
	if m.isTooSmall() {
		t.Fatalf("expected 90x25 to be valid size")
	}
	normalView := m.View()
	if strings.Contains(normalView, "Terminal Window Too Small") {
		t.Fatalf("expected normal view after enlarging, got: %s", normalView)
	}
}

func TestModalAndThemeBordersConsistency(t *testing.T) {
	for _, th := range theme.List() {
		tableBorder := th.TableContainerStyle().GetBorderLeftForeground()
		themeBorder := th.ThemeBoxStyle().GetBorderLeftForeground()
		modalBorder := th.ModalStyle().GetBorderLeftForeground()

		if tableBorder != themeBorder {
			t.Errorf("theme %s: Table border (%v) does not match ThemeBox border (%v)",
				th.ID, tableBorder, themeBorder)
		}
		if tableBorder != modalBorder {
			t.Errorf("theme %s: Table border (%v) does not match Modal border (%v)",
				th.ID, tableBorder, modalBorder)
		}
	}
}

func TestComputeTableColumnsWidth(t *testing.T) {
	for targetWidth := 60; targetWidth <= 160; targetWidth += 10 {
		cols := computeTableColumns(targetWidth)
		tb := table.New(table.WithColumns(cols), table.WithRows([]table.Row{
			{"8080", "1234", "node", "tcp", "LISTEN"},
		}))
		tv := theme.Default().TableContainerStyle().Render(tb.View())
		renderedWidth := lipgloss.Width(tv)
		if renderedWidth != targetWidth {
			t.Fatalf("target %d got rendered %d", targetWidth, renderedWidth)
		}
	}
}

func TestHelpCommandAndKeybindings(t *testing.T) {
	m := newTestModel()
	m.width = 120
	m.height = 30
	m.updateTableDimensions()

	// 1. Open help via '?' key
	mModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = mModel.(Model)
	if !m.isHelp {
		t.Fatalf("expected isHelp to be true after pressing '?'")
	}

	view := m.View()
	if !strings.Contains(view, "Help & Command Reference") {
		t.Fatalf("expected view to contain 'Help & Command Reference', got:\n%s", view)
	}
	if !strings.Contains(view, "KEYBINDINGS") || !strings.Contains(view, "COMMAND ROW") {
		t.Fatalf("expected view to contain keybindings and command row sections, got:\n%s", view)
	}
	if !strings.Contains(view, ":sn") || !strings.Contains(view, ":spn") || !strings.Contains(view, ":theme") {
		t.Fatalf("expected view to describe available commands, got:\n%s", view)
	}

	// 2. Close help via 'Esc'
	mModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = mModel.(Model)
	if m.isHelp {
		t.Fatalf("expected isHelp to be false after pressing Esc")
	}

	// 3. Open help via 'h' key
	mModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = mModel.(Model)
	if !m.isHelp {
		t.Fatalf("expected isHelp to be true after pressing 'h'")
	}

	// Close help via 'q'
	mModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = mModel.(Model)
	if m.isHelp {
		t.Fatalf("expected isHelp to be false after pressing 'q'")
	}
	if cmd != nil {
		t.Fatalf("expected 'q' in help modal to close help without quitting program")
	}

	// 4. Open help via command ':help'
	m.applyCommand("help")
	if !m.isHelp {
		t.Fatalf("expected isHelp to be true after running ':help'")
	}
	m.isHelp = false

	// Open help via command ':h'
	m.applyCommand("h")
	if !m.isHelp {
		t.Fatalf("expected isHelp to be true after running ':h'")
	}
	m.isHelp = false

	// Open help via command ':?'
	m.applyCommand("?")
	if !m.isHelp {
		t.Fatalf("expected isHelp to be true after running ':?'")
	}

	// Close help via 'enter'
	mModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mModel.(Model)
	if m.isHelp {
		t.Fatalf("expected isHelp to be false after pressing Enter")
	}

	// 5. Quit command ':q'
	quitCmd := m.applyCommand("q")
	if quitCmd == nil {
		t.Fatalf("expected ':q' to return tea.Quit cmd")
	}
}

