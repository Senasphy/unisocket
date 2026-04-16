package app

import (
	"fmt"
	"github.com/Senasphy/unisocket/internal/connections"
	"github.com/Senasphy/unisocket/internal/snapshot"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sort"
	"strconv"
	"strings"
)

var (
	pageStyle = lipgloss.NewStyle().
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("31")).
			Padding(0, 1)

	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("222"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	modalStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("69")).
			Padding(1, 2).
			Width(60)

	commandRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	if m.isNamingFile {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				filename := m.filename.Value()
				if filename == "" {
					filename = "portList.json"
				}
				if !strings.HasSuffix(strings.ToLower(filename), ".json") {
					filename += ".json"
				}
				m.isNamingFile = false
				return m, snapshot.SaveStateToFile(m.table.Rows(), filename)
			case "esc":
				m.isNamingFile = false
				return m, nil
			}
		}
		m.filename, cmd = m.filename.Update(msg)
		return m, cmd
	}

	if m.isSearching {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.isSearching = false
				return m, nil
			case "esc":
				m.isSearching = false
				m.searchQuery = ""
				m.applyRows()
				return m, nil
			}
		}
		m.searchInput, cmd = m.searchInput.Update(msg)
		m.searchQuery = strings.TrimSpace(m.searchInput.Value())
		m.applyRows()
		return m, cmd
	}

	if m.isCommandMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				command := strings.TrimSpace(m.commandInput.Value())
				command = strings.TrimPrefix(command, ":")
				m.commandInput.Reset()
				m.isCommandMode = false
				m.applyCommand(command)
				m.applyRows()
				return m, nil
			case "esc":
				m.commandInput.Reset()
				m.isCommandMode = false
				return m, nil
			}
		}
		m.commandInput, cmd = m.commandInput.Update(msg)
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		tableHeight := msg.Height - 12
		tableHeight = max(5,tableHeight)
		m.table.SetHeight(tableHeight)
		return m, nil
	case tickMsg:
		return m, tea.Batch(m.fetchDataCmd(), tick())
	case DataArrivalMsg:
		m.allRows = msg
		m.applyRows()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "s":
			m.isNamingFile = true
			m.filename = textinput.New()
			m.filename.Placeholder = "filename.json"
			m.filename.Focus()
			return m, nil
		case "/":
			m.isSearching = true
			m.searchInput = textinput.New()
			m.searchInput.Placeholder = "search connections..."
			m.searchInput.SetValue(m.searchQuery)
			m.searchInput.Focus()
			return m, nil
		case ":":
			m.isCommandMode = true
			m.commandInput = textinput.New()
			m.commandInput.Placeholder = "sn | spn | spid"
			m.commandInput.Focus()
			return m, nil
		case "K":
			return m, m.killSelectedProcess()
		}
	case snapshot.SavedMessageString:
		m.statusMessage = string(msg)
		return m, nil

	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	title := titleStyle.Render("UniSocket  •  Live Connection Monitor")
	tableView := tableStyle.Render(m.table.View())
	content := title + "\n" + tableView + "\n"

	if m.statusMessage != "" {
		content += "\n" + statusStyle.Render(m.statusMessage) + "\n"
	}

	searchHint := "none"
	if m.searchQuery != "" {
		searchHint = m.searchQuery
	}
	sortHint := "none"
	if m.sortMode != "" {
		sortHint = m.sortMode
	}
	footer := fmt.Sprintf("↑/↓: Navigate • /: Search • K: Kill process • s: Snapshot • : Command row • q: Quit   |   search=%s sort=%s", searchHint, sortHint)
	content += footerStyle.Render(footer)

	if m.isCommandMode {
		content += "\n\n" + commandRowStyle.Render(":"+m.commandInput.View())
	}

	if m.isSearching {
		content += "\n\n" + commandRowStyle.Render("/"+m.searchInput.View())
	}

	view := pageStyle.Render(content)
	if m.isNamingFile {
		modal := modalStyle.Render(
			fmt.Sprintf("Save snapshot as:\n\n%s\n\nEnter: Save    Esc: Cancel", m.filename.View()),
		)
		if m.width > 0 {
			return view + "\n\n" + lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(modal)
		}
		return view + "\n\n" + modal
	}

	return view
}

func (m *Model) applyRows() {
	rows := make([]table.Row, len(m.allRows))
	for i, row := range m.allRows {
		rows[i] = append(table.Row(nil), row...)
	}

	if m.sortMode != "" {
		sort.SliceStable(rows, func(i, j int) bool {
			switch m.sortMode {
			case "name":
				return strings.ToLower(rows[i][0]) < strings.ToLower(rows[j][0])
			case "process-name":
				return strings.ToLower(rows[i][2]) < strings.ToLower(rows[j][2])
			case "pid":
				left, errLeft := strconv.Atoi(rows[i][1])
				right, errRight := strconv.Atoi(rows[j][1])
				if errLeft != nil && errRight != nil {
					return rows[i][1] < rows[j][1]
				}
				if errLeft != nil {
					return false
				}
				if errRight != nil {
					return true
				}
				return left < right
			default:
				return false
			}
		})
	}

	filtered := rows
	query := strings.ToLower(strings.TrimSpace(m.searchQuery))
	if query != "" {
		filtered = make([]table.Row, 0, len(rows))
		for _, row := range rows {
			joined := strings.ToLower(strings.Join(row, " "))
			if strings.Contains(joined, query) {
				filtered = append(filtered, row)
			}
		}
	}

	m.table.SetRows(filtered)
	if len(filtered) == 0 {
		m.table.SetCursor(0)
		return
	}
	if m.table.Cursor() >= len(filtered) {
		m.table.SetCursor(len(filtered) - 1)
	}
}

func (m *Model) applyCommand(command string) {
	switch command {
	case "sn":
		m.sortMode = "name"
		m.statusMessage = "Sorted by name"
	case "spn":
		m.sortMode = "process-name"
		m.statusMessage = "Sorted by process name"
	case "spid":
		m.sortMode = "pid"
		m.statusMessage = "Sorted by PID"
	case "":
		m.statusMessage = "Command cancelled"
	default:
		m.statusMessage = fmt.Sprintf("Unknown command: %s", command)
	}
}

func (m *Model) killSelectedProcess() tea.Cmd {
	row := m.table.SelectedRow()
	if len(row) < 2 {
		m.statusMessage = "No row selected"
		return nil
	}
	pid := strings.TrimSpace(row[1])
	if pid == "" || pid == "-" {
		m.statusMessage = "Selected row has no killable PID"
		return nil
	}
	if err := connections.KillProcess(pid); err != nil {
		m.statusMessage = fmt.Sprintf("Failed to kill PID %s: %v", pid, err)
		return nil
	}
	m.statusMessage = fmt.Sprintf("Killed PID %s", pid)
	return m.fetchDataCmd()
}
