package app

import (
	"os"
	"time"

	"github.com/Senasphy/unisocket/internal/config"
	"github.com/Senasphy/unisocket/internal/connections"
	"github.com/Senasphy/unisocket/internal/theme"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

type tickMsg time.Time
type DataArrivalMsg []table.Row

type Model struct {
	table            table.Model
	statusMessage    string
	filename         textinput.Model
	searchInput      textinput.Model
	commandInput     textinput.Model
	isNamingFile     bool
	isKillConfirm    bool
	isSearching      bool
	isCommandMode    bool
	isThemePicker    bool
	isHelp           bool
	themePickerIndex int
	previousThemeID  string
	currentTheme     theme.Theme
	cfg              *config.Config
	filters          connections.Filters
	allRows          []table.Row
	searchQuery      string
	sortMode         string
	pendingKillPID   string
	pendingKillName  string
	width            int
	height           int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.fetchDataCmd(), m.tickCmd())
}

func New(t table.Model, filter connections.Filters, cfg *config.Config, th ...theme.Theme) tea.Model {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	activeTheme := theme.Default()
	if len(th) > 0 {
		activeTheme = th[0]
	} else if cfg.Theme != "" {
		if found, ok := theme.Get(cfg.Theme); ok {
			activeTheme = found
		}
	}
	t.SetStyles(activeTheme.TableStyles())

	sortMode := ""
	if cfg.DefaultSort != "" && cfg.DefaultSort != "none" {
		sortMode = cfg.DefaultSort
	}

	m := Model{
		table:        t,
		filters:      filter,
		allRows:      t.Rows(),
		currentTheme: activeTheme,
		cfg:          cfg,
		sortMode:     sortMode,
	}
	if sortMode != "" {
		m.applyRows()
	}
	if w, h, err := term.GetSize(os.Stdout.Fd()); err == nil && w > 0 && h > 0 {
		m.width = w
		m.height = h
		m.updateTableDimensions()
	}
	return m
}

func (m *Model) setTheme(th theme.Theme) {
	m.currentTheme = th
	m.table.SetStyles(th.TableStyles())
}

func (m *Model) saveTheme(themeID string) {
	if m.cfg != nil {
		m.cfg.Theme = themeID
		_ = m.cfg.Save()
	} else {
		_ = config.SaveTheme(themeID)
	}
}

func (m *Model) openThemePicker() {
	m.isThemePicker = true
	m.previousThemeID = m.currentTheme.ID
	themes := theme.List()
	m.themePickerIndex = 0
	for i, th := range themes {
		if th.ID == m.currentTheme.ID {
			m.themePickerIndex = i
			break
		}
	}
	m.updateTableDimensions()
}

func (m Model) tickCmd() tea.Cmd {
	interval := 2
	if m.cfg != nil && m.cfg.RefreshInterval > 0 {
		interval = m.cfg.RefreshInterval
	}
	return tea.Tick(time.Duration(interval)*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) fetchDataCmd() tea.Cmd {
	return func() tea.Msg {
		rows := connections.GetConnections(m.filters)
		return DataArrivalMsg(rows)
	}
}
