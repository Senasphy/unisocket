package app

import (
	"github.com/Senasphy/unisocket/internal/connections"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type tickMsg time.Time
type DataArrivalMsg []table.Row

type Model struct {
	table         table.Model
	statusMessage string
	filename      textinput.Model
	searchInput   textinput.Model
	commandInput  textinput.Model
	isNamingFile  bool
	isSearching   bool
	isCommandMode bool
	filters       connections.Filters
	allRows       []table.Row
	searchQuery   string
	sortMode      string
	width         int
	height        int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.fetchDataCmd(), tick())
}

func New(t table.Model, filter connections.Filters) tea.Model {
	return Model{
		table:   t,
		filters: filter,
		allRows: t.Rows(),
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) fetchDataCmd() tea.Cmd {
	return func() tea.Msg {
		rows := connections.GetConnections(m.filters)
		return DataArrivalMsg(rows)
	}
}
