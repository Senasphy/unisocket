package app

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Senasphy/unisocket/internal/connections"
	"github.com/Senasphy/unisocket/internal/snapshot"
	"github.com/Senasphy/unisocket/internal/theme"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

const (
	minTerminalWidth  = 64
	minTerminalHeight = 12
	themeBoxWidth     = 38
	gapWidth          = 2
	pagePaddingH      = 4 // pageStyle padding(1, 2)
)

var (
	pageStyle = lipgloss.NewStyle().
			Padding(1, 2)
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.isTooSmall() {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			m.updateTableDimensions()
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return m, tea.Quit
			}
		}
		return m, nil
	}

	var cmd tea.Cmd
	if m.isHelp {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			m.updateTableDimensions()
			return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "enter", "q", "?", "h", "space":
				m.isHelp = false
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
		}
		return m, nil
	}

	if m.isNamingFile {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				defaultName := "portList.json"
				if m.cfg != nil && m.cfg.SnapshotFilename != "" {
					defaultName = m.cfg.SnapshotFilename
				}
				filename := m.filename.Value()
				if filename == "" {
					filename = defaultName
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

	if m.isKillConfirm {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			m.isKillConfirm = false
			if keyMsg.String() == "enter" {
				return m, m.terminatePendingProcess()
			}
			m.pendingKillPID = ""
			m.pendingKillName = ""
			m.statusMessage = "Process termination cancelled"
			return m, nil
		}
		return m, nil
	}

	if m.isSearching {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.isSearching = false
				m.searchQuery = strings.TrimSpace(m.searchInput.Value())
				m.applyRows()
				if m.searchQuery != "" {
					m.statusMessage = fmt.Sprintf("Search: %s (Esc to clear)", m.searchQuery)
				}
				return m, nil
			case "esc", "ctrl+c":
				m.isSearching = false
				m.searchQuery = ""
				m.searchInput.Reset()
				m.applyRows()
				m.statusMessage = "Search cleared"
				return m, nil
			case "down", "tab":
				m.isSearching = false
				m.searchQuery = strings.TrimSpace(m.searchInput.Value())
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
				cmd := m.applyCommand(command)
				m.applyRows()
				return m, cmd
			case "esc", "ctrl+c":
				m.commandInput.Reset()
				m.isCommandMode = false
				m.statusMessage = "Command cancelled"
				return m, nil
			}
		}
		m.commandInput, cmd = m.commandInput.Update(msg)
		return m, cmd
	}

	if m.isThemePicker {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			themes := theme.List()
			switch msg.String() {
			case "up", "k":
				if m.themePickerIndex > 0 {
					m.themePickerIndex--
				} else {
					m.themePickerIndex = len(themes) - 1
				}
				m.setTheme(themes[m.themePickerIndex])
				return m, nil
			case "down", "j":
				if m.themePickerIndex < len(themes)-1 {
					m.themePickerIndex++
				} else {
					m.themePickerIndex = 0
				}
				m.setTheme(themes[m.themePickerIndex])
				return m, nil
			case "enter":
				m.isThemePicker = false
				selected := themes[m.themePickerIndex]
				m.setTheme(selected)
				m.saveTheme(selected.ID)
				m.statusMessage = fmt.Sprintf("Theme set to %s", selected.Name)
				m.updateTableDimensions()
				return m, nil
			case "esc", "q":
				m.isThemePicker = false
				if prev, ok := theme.Get(m.previousThemeID); ok {
					m.setTheme(prev)
				}
				m.statusMessage = "Theme change cancelled"
				m.updateTableDimensions()
				return m, nil
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateTableDimensions()
		return m, nil
	case tickMsg:
		return m, tea.Batch(m.fetchDataCmd(), m.tickCmd())
	case DataArrivalMsg:
		m.allRows = msg
		m.applyRows()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.searchQuery != "" {
				m.searchQuery = ""
				m.searchInput.Reset()
				m.applyRows()
				m.statusMessage = "Search cleared"
				return m, nil
			}
			return m, tea.Quit
		case "ctrl+c", "q":
			return m, tea.Quit
		case "t":
			m.openThemePicker()
			return m, nil
		case "s":
			m.isNamingFile = true
			m.filename = textinput.New()
			defaultName := "portList.json"
			if m.cfg != nil && m.cfg.SnapshotFilename != "" {
				defaultName = m.cfg.SnapshotFilename
			}
			m.filename.Placeholder = defaultName
			m.filename.Focus()
			return m, nil
		case "/":
			m.isSearching = true
			m.searchInput = textinput.New()
			m.searchInput.Placeholder = "search connections..."
			m.searchInput.SetValue(m.searchQuery)
			m.searchInput.Focus()
			return m, nil
		case "?", "h":
			m.isHelp = true
			m.statusMessage = "Help opened (press Esc to close)"
			return m, nil
		case ":":
			m.isCommandMode = true
			m.commandInput = textinput.New()
			m.commandInput.Placeholder = "help | sn | spn | spid | theme [name] | q"
			m.commandInput.Focus()
			return m, nil
		case "K":
			m.startKillConfirmation()
			return m, nil
		}
	case snapshot.SavedMessageString:
		m.statusMessage = string(msg)
		return m, nil

	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.isTooSmall() {
		return m.renderTooSmallView()
	}

	title := m.currentTheme.TitleStyle().Render("UniSocket  •  Live Connection Monitor")
	tableView := m.currentTheme.TableContainerStyle().Render(m.table.View())

	var mainSection string
	sideBySide := false
	if m.isThemePicker {
		tableWidth := lipgloss.Width(tableView)
		tableHeight := lipgloss.Height(tableView)
		sideBySide = m.isSideBySide(tableWidth)
		if sideBySide {
			naturalThemeBox := m.renderThemePickerBox(tableView, true, 0)
			themeHeight := lipgloss.Height(naturalThemeBox)
			matchHeight := max(tableHeight, themeHeight)

			var themeBox string
			if themeHeight != matchHeight {
				themeBox = m.renderThemePickerBox(tableView, true, matchHeight)
			} else {
				themeBox = naturalThemeBox
			}

			if tableHeight != matchHeight {
				tableView = m.currentTheme.TableContainerStyle().
					Height(matchHeight - 2).
					Render(m.table.View())
			}

			mainSection = lipgloss.JoinHorizontal(lipgloss.Top, tableView, "  ", themeBox)
		} else {
			themeBox := m.renderThemePickerBox(tableView, false, 0)
			mainSection = lipgloss.JoinVertical(lipgloss.Left, tableView, "\n"+themeBox)
		}
	} else {
		mainSection = tableView
	}

	totalWidth := lipgloss.Width(mainSection)
	tableWidth := lipgloss.Width(tableView)

	content := title + "\n" + mainSection

	if m.isSearching {
		content += "\n" + m.renderSearchBox(totalWidth)
	}

	if m.isCommandMode {
		content += "\n" + m.renderCommandBox(totalWidth)
	}

	bottomDeck := m.renderBottomDeck(totalWidth, tableWidth, sideBySide)
	content += "\n" + bottomDeck

	view := pageStyle.Render(content)
	if m.isNamingFile {
		mWidth := m.modalWidthForTable(tableView)
		mStyle := m.currentTheme.ModalStyle()
		innerW := mWidth - mStyle.GetHorizontalBorderSize()
		if innerW < 20 {
			innerW = 20
		}
		modal := mStyle.Width(innerW).Render(
			fmt.Sprintf("Save snapshot as:\n\n%s\n\nEnter: Save    Esc: Cancel", m.filename.View()),
		)
		return view + "\n\n" + m.centerWithinTable(tableView, modal)
	}

	if m.isKillConfirm {
		mWidth := m.modalWidthForTable(tableView)
		mStyle := m.currentTheme.ModalStyle()
		innerW := mWidth - mStyle.GetHorizontalBorderSize()
		if innerW < 20 {
			innerW = 20
		}
		modal := mStyle.Width(innerW).Render(
			fmt.Sprintf(
				"Confirm Termination\n\nPID: %s   |   Process: %s\n\n[Enter] Terminate    [Any other key] Cancel",
				m.pendingKillPID,
				m.pendingKillName,
			),
		)
		return view + "\n\n" + m.centerWithinTable(tableView, modal)
	}

	if m.isHelp {
		modal := m.renderHelpModal(tableView)
		return pageStyle.Render(title + "\n" + m.centerWithinTable(tableView, modal))
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

func (m *Model) applyCommand(command string) tea.Cmd {
	trimmed := strings.TrimSpace(command)
	lower := strings.ToLower(trimmed)
	switch {
	case lower == "help" || lower == "h" || lower == "?":
		m.isHelp = true
		m.statusMessage = "Help opened (press Esc to close)"
		return nil
	case lower == "q" || lower == "quit":
		return tea.Quit
	case lower == "sn":
		m.sortMode = "name"
		m.statusMessage = "Sorted by name"
		return nil
	case lower == "spn":
		m.sortMode = "process-name"
		m.statusMessage = "Sorted by process name"
		return nil
	case lower == "spid":
		m.sortMode = "pid"
		m.statusMessage = "Sorted by PID"
		return nil
	case lower == "theme" || lower == "t":
		m.openThemePicker()
		return nil
	case strings.HasPrefix(lower, "theme ") || strings.HasPrefix(lower, "t "):
		parts := strings.SplitN(trimmed, " ", 2)
		if len(parts) == 2 {
			themeQuery := strings.TrimSpace(parts[1])
			if th, ok := theme.Get(themeQuery); ok {
				m.setTheme(th)
				m.saveTheme(th.ID)
				m.statusMessage = fmt.Sprintf("Theme set to %s", th.Name)
			} else {
				m.statusMessage = fmt.Sprintf("Unknown theme '%s'. Press 't' or run ':theme' to browse", themeQuery)
			}
		}
		return nil
	case lower == "":
		m.statusMessage = "Command cancelled"
		return nil
	default:
		m.statusMessage = fmt.Sprintf("Unknown command: %s (type :help for list)", command)
		return nil
	}
}

func (m Model) terminalWidth() int {
	if m.width > 0 {
		return m.width
	}
	if w, _, err := term.GetSize(os.Stdout.Fd()); err == nil && w > 0 {
		return w
	}
	return 0
}

func (m Model) isSideBySideForWidth(termW int) bool {
	if termW <= 0 {
		return false
	}
	return termW >= 60+themeBoxWidth+gapWidth+pagePaddingH
}

func (m Model) isSideBySide(tableWidth int) bool {
	termW := m.terminalWidth()
	return m.isSideBySideForWidth(termW)
}

func (m Model) isTooSmall() bool {
	return m.width > 0 && m.height > 0 && (m.width < minTerminalWidth || m.height < minTerminalHeight)
}

func (m Model) renderTooSmallView() string {
	boxStyle := m.currentTheme.ThemeBoxStyle().Padding(1, 2)
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.currentTheme.Highlight).
		Render("Terminal Window Too Small")

	msg := fmt.Sprintf(
		"UniSocket requires a minimum terminal size of %dx%d.\n\n"+
			"Current size: %dx%d\n\n"+
			"Please resize or enlarge your terminal window.\n"+
			"Press 'q' or Ctrl+C to quit.",
		minTerminalWidth, minTerminalHeight,
		m.width, m.height,
	)
	text := lipgloss.NewStyle().Foreground(m.currentTheme.Text).Render(msg)
	box := boxStyle.Render(title + "\n\n" + text)

	if m.width >= lipgloss.Width(box) && m.height >= lipgloss.Height(box) {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}
	return box
}

func computeTableColumns(targetContainerWidth int) []table.Column {
	colSum := targetContainerWidth - 14
	if colSum < 46 {
		colSum = 46
	}

	pidW := 9
	protoW := 7
	stateW := 12
	if colSum >= 80 {
		pidW = 11
		protoW = 8
		stateW = 13
	}
	if colSum >= 110 {
		pidW = 12
		protoW = 10
		stateW = 14
	}

	fixedSum := pidW + protoW + stateW
	rem := colSum - fixedSum
	if rem < 18 {
		rem = 18
	}

	portW := (rem * 45) / 100
	if portW < 8 {
		portW = 8
	}
	pnameW := rem - portW

	return []table.Column{
		{Title: "PORT", Width: portW},
		{Title: "PROCESS ID", Width: pidW},
		{Title: "PROCESS NAME", Width: pnameW},
		{Title: "PROTOCOL", Width: protoW},
		{Title: "STATE", Width: stateW},
	}
}

func (m *Model) updateTableDimensions() {
	termW := m.terminalWidth()
	if termW <= 0 {
		return
	}

	if m.height > 0 {
		tableHeight := m.height - 12
		if tableHeight < 5 {
			tableHeight = 5
		}
		m.table.SetHeight(tableHeight)
	}

	availW := termW - pagePaddingH
	if availW < 50 {
		availW = 50
	}

	targetTableWidth := availW
	if m.isThemePicker && m.isSideBySideForWidth(termW) {
		targetTableWidth = availW - themeBoxWidth - gapWidth
	}

	if targetTableWidth < 50 {
		targetTableWidth = 50
	}

	m.table.SetColumns(computeTableColumns(targetTableWidth))
}

func (m Model) renderThemePickerBox(tableView string, isSideBySide bool, targetHeight int) string {
	themes := theme.List()

	cursorStyle := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Highlight)
	selectedStyle := lipgloss.NewStyle().Foreground(m.currentTheme.Status)
	itemStyle := lipgloss.NewStyle().Foreground(m.currentTheme.Text)
	descStyle := lipgloss.NewStyle().Foreground(m.currentTheme.Muted)

	items := make([]string, len(themes))
	for i, th := range themes {
		cursor := "  "
		if i == m.themePickerIndex {
			cursor = "▸ "
		}

		indicator := "  "
		if th.ID == m.previousThemeID {
			indicator = "✓ "
		}

		nameCol := fmt.Sprintf("%-26s", th.Name)
		if isSideBySide {
			if i == m.themePickerIndex {
				items[i] = cursorStyle.Render(cursor) + selectedStyle.Render(indicator) + cursorStyle.Render(nameCol)
			} else {
				items[i] = itemStyle.Render(cursor) + selectedStyle.Render(indicator) + itemStyle.Render(nameCol)
			}
		} else {
			if i == m.themePickerIndex {
				items[i] = cursorStyle.Render(cursor) + selectedStyle.Render(indicator) + cursorStyle.Render(nameCol) + " " + descStyle.Render(th.Description)
			} else {
				items[i] = itemStyle.Render(cursor) + selectedStyle.Render(indicator) + itemStyle.Render(nameCol) + " " + descStyle.Render(th.Description)
			}
		}
	}

	var title string
	var controls string
	boxStyle := m.currentTheme.ThemeBoxStyle()

	if isSideBySide {
		title = lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Highlight).Render("Themes (Preview)")
		controls = m.currentTheme.FooterStyle().Render("↑/↓ or j/k: Preview\nEnter: Apply • Esc: Cancel")
		body := title + "\n\n" + strings.Join(items, "\n") + "\n\n" + controls

		innerW := themeBoxWidth - boxStyle.GetHorizontalBorderSize()
		if innerW < 10 {
			innerW = 10
		}
		boxStyle = boxStyle.Width(innerW)
		if targetHeight > 2 {
			boxStyle = boxStyle.Height(targetHeight - 2)
		}
		return boxStyle.Render(body)
	}

	title = lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Highlight).Render("Theme Switcher (Live Preview)")
	controls = m.currentTheme.FooterStyle().Render("↑/↓ or j/k: Preview   •   Enter: Apply   •   Esc: Cancel")
	body := title + "\n\n" + strings.Join(items, "\n") + "\n\n" + controls

	modalWidth := m.modalWidthForTable(tableView)
	if modalWidth < 68 {
		modalWidth = 68
	}
	innerW := modalWidth - boxStyle.GetHorizontalBorderSize()
	if innerW < 10 {
		innerW = 10
	}
	return boxStyle.Width(innerW).Render(body)
}

func (m Model) renderHelpModal(tableView string) string {
	tableWidth := lipgloss.Width(tableView)
	modalWidth := m.modalWidthForTable(tableView)
	if modalWidth < 68 && tableWidth >= 68 {
		modalWidth = 68
	}
	if modalWidth > 74 {
		modalWidth = 74
	}

	mStyle := m.currentTheme.ModalStyle()
	innerW := modalWidth - mStyle.GetHorizontalBorderSize()
	if innerW < 40 {
		innerW = 40
	}
	mStyle = mStyle.Width(innerW)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.currentTheme.Highlight).
		Render("UniSocket  •  Help & Command Reference")

	sectionHeader := func(s string) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(m.currentTheme.Highlight).
			Render(s)
	}

	badge := func(k string) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(m.currentTheme.TitleText).
			Background(m.currentTheme.Highlight).
			Padding(0, 1).
			Render(k)
	}

	cmdLabel := func(c string) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(m.currentTheme.Status).
			Render(c)
	}

	desc := func(d string) string {
		return lipgloss.NewStyle().
			Foreground(m.currentTheme.Text).
			Render(d)
	}

	subDesc := func(d string) string {
		return lipgloss.NewStyle().
			Foreground(m.currentTheme.Muted).
			Render(d)
	}

	// 2-column layout for wide modals
	if innerW >= 62 {
		colW := (innerW - 4) / 2
		formatColRow := func(left, right string) string {
			padLeft := left
			lw := lipgloss.Width(left)
			if lw < colW {
				padLeft += strings.Repeat(" ", colW-lw)
			}
			return padLeft + "  " + right
		}

		kRow := func(k, d string) string {
			b := badge(k)
			bw := lipgloss.Width(b)
			sp := " "
			if bw < 10 {
				sp = strings.Repeat(" ", 10-bw)
			}
			return b + sp + desc(d)
		}

		cRow := func(c, d string) string {
			cl := cmdLabel(c)
			cw := lipgloss.Width(cl)
			sp := " "
			if cw < 14 {
				sp = strings.Repeat(" ", 14-cw)
			}
			return cl + sp + desc(d)
		}

		rows := []string{
			formatColRow(sectionHeader("KEYBINDINGS"), sectionHeader("COMMAND ROW (:)")),
			formatColRow(kRow("↑/↓", "Navigate table"), cRow(":help, :h", "Show this guide")),
			formatColRow(kRow("/", "Search connections"), cRow(":sn", "Sort by port name")),
			formatColRow(kRow("t", "Theme switcher"), cRow(":spn", "Sort by process")),
			formatColRow(kRow("s", "Snapshot to JSON"), cRow(":spid", "Sort by PID")),
			formatColRow(kRow("K", "Terminate process"), cRow(":theme [th]", "Switch theme")),
			formatColRow(kRow("?, h", "Toggle help guide"), cRow(":q, :quit", "Quit UniSocket")),
			formatColRow(kRow("Esc", "Clear / dismiss"), ""),
			formatColRow(kRow("q", "Quit program"), ""),
		}

		searchHint := subDesc("Search Mode (/): Esc / Ctrl+C to cancel • Tab / ↓ to move to table")
		footer := m.currentTheme.FooterStyle().Render("Press [Esc], [Enter], [?], or [q] to close help")

		body := title + "\n\n" + strings.Join(rows, "\n") + "\n\n" + searchHint + "\n\n" + footer
		return mStyle.Render(body)
	}

	// Compact single-column layout for narrower terminals
	row := func(k, d string) string {
		b := badge(k)
		bw := lipgloss.Width(b)
		sp := " "
		if bw < 12 {
			sp = strings.Repeat(" ", 12-bw)
		}
		return b + sp + desc(d)
	}

	lines := []string{
		title,
		"",
		sectionHeader("KEYBINDINGS & ACTIONS"),
		row("↑/↓, j/k", "Navigate connections"),
		row("/", "Search & filter rows"),
		row("t", "Switch color theme"),
		row("s", "Save JSON snapshot"),
		row("K", "Terminate selected PID"),
		row("?, h", "Open this help guide"),
		row("Esc", "Clear filter / dismiss"),
		row("q", "Quit UniSocket"),
		"",
		sectionHeader("COMMANDS (Press :)"),
		row(":help, :h", "Open this reference"),
		row(":sn/:spn", "Sort: Port / Process"),
		row(":spid", "Sort by Process ID"),
		row(":theme [th]", "Change color theme"),
		row(":q", "Quit UniSocket"),
		"",
		m.currentTheme.FooterStyle().Render("Press [Esc], [Enter], [?], or [q] to close"),
	}

	return mStyle.Render(strings.Join(lines, "\n"))
}

func (m *Model) startKillConfirmation() {
	row := m.table.SelectedRow()
	if len(row) < 2 {
		m.statusMessage = "No row selected"
		return
	}
	pid := strings.TrimSpace(row[1])
	if pid == "" || pid == "-" {
		m.statusMessage = "Selected row has no killable PID"
		return
	}
	processName := "Unknown"
	if len(row) > 2 && strings.TrimSpace(row[2]) != "" {
		processName = strings.TrimSpace(row[2])
	}
	m.pendingKillPID = pid
	m.pendingKillName = processName
	m.isKillConfirm = true
}

func (m *Model) terminatePendingProcess() tea.Cmd {
	pid := strings.TrimSpace(m.pendingKillPID)
	processName := strings.TrimSpace(m.pendingKillName)
	m.pendingKillPID = ""
	m.pendingKillName = ""
	if pid == "" {
		m.statusMessage = "No process selected"
		return nil
	}
	if err := connections.KillProcess(pid); err != nil {
		m.statusMessage = fmt.Sprintf("Failed to terminate PID %s: %v", pid, err)
		return nil
	}
	m.statusMessage = fmt.Sprintf("Terminated PID %s (%s)", pid, processName)
	return m.fetchDataCmd()
}

func (m Model) centerWithinTable(tableView string, block string) string {
	tableWidth := lipgloss.Width(tableView)
	blockWidth := lipgloss.Width(block)
	if tableWidth <= blockWidth || tableWidth == 0 {
		return block
	}
	return lipgloss.NewStyle().Width(tableWidth).Align(lipgloss.Center).Render(block)
}

func (m Model) modalWidthForTable(tableView string) int {
	tableWidth := lipgloss.Width(tableView)
	if tableWidth <= 0 {
		return 58
	}
	width := tableWidth - 4
	if width < 42 {
		return 42
	}
	if width > 76 {
		return 76
	}
	return width
}

func (m Model) renderSearchBox(width int) string {
	boxStyle := m.currentTheme.ThemeBoxStyle()
	innerW := width - boxStyle.GetHorizontalBorderSize()
	if innerW < 10 {
		innerW = 10
	}
	boxStyle = boxStyle.Width(innerW)

	title := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Highlight).Render("Search Connections")
	prompt := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Highlight).Render("/ ")
	inputView := prompt + m.searchInput.View()
	hints := m.currentTheme.FooterStyle().Render("Esc / Ctrl+C: Cancel & Exit  •  Enter: Keep Filter  •  Tab / ↓: Move to Table")

	body := title + "\n" + inputView + "\n" + hints
	return boxStyle.Render(body)
}

func (m Model) renderCommandBox(width int) string {
	boxStyle := m.currentTheme.ThemeBoxStyle()
	innerW := width - boxStyle.GetHorizontalBorderSize()
	if innerW < 10 {
		innerW = 10
	}
	boxStyle = boxStyle.Width(innerW)

	title := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Highlight).Render("Command Row")
	prompt := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Highlight).Render(": ")
	inputView := prompt + m.commandInput.View()
	hints := m.currentTheme.FooterStyle().Render(":help  •  :sn/:spn/:spid  •  :theme [name]  •  :q  •  Esc: Cancel")

	body := title + "\n" + inputView + "\n" + hints
	return boxStyle.Render(body)
}

func (m Model) renderBottomDeck(totalWidth int, tableWidth int, isSideBySide bool) string {
	if isSideBySide {
		leftWidth := tableWidth
		rightWidth := totalWidth - leftWidth - gapWidth
		if rightWidth < 28 {
			rightWidth = 28
		}
		leftCard := m.renderControlsCard(leftWidth)
		rightCard := m.renderMonitorCard(rightWidth)
		return lipgloss.JoinHorizontal(lipgloss.Top, leftCard, "  ", rightCard)
	}

	if totalWidth >= 76 {
		rightWidth := 30
		leftWidth := totalWidth - rightWidth - gapWidth
		leftCard := m.renderControlsCard(leftWidth)
		rightCard := m.renderMonitorCard(rightWidth)
		return lipgloss.JoinHorizontal(lipgloss.Top, leftCard, "  ", rightCard)
	}

	return m.renderUnifiedCard(totalWidth)
}

func (m Model) renderControlsCard(width int) string {
	cardStyle := m.currentTheme.ThemeBoxStyle()
	innerW := width - cardStyle.GetHorizontalBorderSize()
	if innerW < 10 {
		innerW = 10
	}
	cardStyle = cardStyle.Width(innerW).Height(4)

	title := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Highlight).Render("Commands")

	badge := func(k, a string) string {
		kStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(m.currentTheme.TitleText).
			Background(m.currentTheme.Highlight).
			Padding(0, 1).
			Render(k)
		aStyle := lipgloss.NewStyle().
			Foreground(m.currentTheme.Text).
			Render(" " + a)
		return kStyle + aStyle
	}

	var row1, row2, hint string
	if width >= 65 {
		row1 = fmt.Sprintf("%s   %s   %s   %s",
			badge("↑/↓", "Nav"),
			badge("/", "Search"),
			badge("t", "Theme"),
			badge("s", "Snapshot"),
		)
		row2 = fmt.Sprintf("%s   %s   %s   %s",
			badge("?", "Help"),
			badge("K", "Kill PID"),
			badge(":", "Command"),
			badge("q", "Quit"),
		)
		if m.searchQuery != "" {
			row2 += "   " + badge("Esc", "Clear")
		}
		hint = m.currentTheme.FooterStyle().Render("Tip: ?/h or :help for command reference • :sn/:spn/:spid sort")
	} else {
		row1 = fmt.Sprintf("%s  %s  %s",
			badge("↑↓", "Nav"),
			badge("/", "Find"),
			badge("t", "Theme"),
		)
		row2 = fmt.Sprintf("%s  %s  %s  %s",
			badge("?", "Help"),
			badge("s", "Snap"),
			badge("K", "Kill"),
			badge("q", "Quit"),
		)
		if m.searchQuery != "" {
			hint = badge("Esc", "Clear Search")
		} else {
			hint = m.currentTheme.FooterStyle().Render("?: Help • :sn/:spn/:spid sort")
		}
	}

	body := title + "\n" + row1 + "\n" + row2 + "\n" + hint
	return cardStyle.Render(body)
}

func (m Model) renderMonitorCard(width int) string {
	cardStyle := m.currentTheme.ThemeBoxStyle()
	innerW := width - cardStyle.GetHorizontalBorderSize()
	if innerW < 10 {
		innerW = 10
	}
	cardStyle = cardStyle.Width(innerW).Height(4)

	title := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Highlight).Render("Monitor")

	labelStyle := lipgloss.NewStyle().Foreground(m.currentTheme.Muted)
	valStyle := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Text)

	themeName := m.currentTheme.Name
	if len(themeName) > 16 {
		themeName = m.currentTheme.ID
	}
	line1 := labelStyle.Render("Theme: ") + valStyle.Render(themeName)

	sortMode := "none"
	if m.sortMode != "" {
		sortMode = m.sortMode
	}
	line2 := labelStyle.Render("Sort: ") + valStyle.Render(sortMode) +
		labelStyle.Render(" • ") +
		lipgloss.NewStyle().Foreground(m.currentTheme.Status).Render(fmt.Sprintf("%d conns", len(m.allRows)))

	var line3 string
	if m.statusMessage != "" {
		statusText := m.statusMessage
		maxLen := width - 6
		if maxLen > 5 && len(statusText) > maxLen {
			statusText = statusText[:maxLen-3] + "..."
		}
		line3 = lipgloss.NewStyle().Foreground(m.currentTheme.Status).Bold(true).Render("● " + statusText)
	} else if m.searchQuery != "" {
		searchQuery := m.searchQuery
		maxLen := width - 15
		if maxLen > 5 && len(searchQuery) > maxLen {
			searchQuery = searchQuery[:maxLen-3] + "..."
		}
		line3 = labelStyle.Render("Search: ") + lipgloss.NewStyle().Foreground(m.currentTheme.Highlight).Render(searchQuery)
	} else {
		interval := 2
		if m.cfg != nil && m.cfg.RefreshInterval > 0 {
			interval = m.cfg.RefreshInterval
		}
		line3 = labelStyle.Render(fmt.Sprintf("● Live feed (%ds)", interval))
	}

	body := title + "\n" + line1 + "\n" + line2 + "\n" + line3
	return cardStyle.Render(body)
}

func (m Model) renderUnifiedCard(width int) string {
	cardStyle := m.currentTheme.ThemeBoxStyle()
	innerW := width - cardStyle.GetHorizontalBorderSize()
	if innerW < 10 {
		innerW = 10
	}
	cardStyle = cardStyle.Width(innerW)

	title := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.Highlight).Render("Commands & Status")

	badge := func(k, a string) string {
		kStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(m.currentTheme.TitleText).
			Background(m.currentTheme.Highlight).
			Padding(0, 1).
			Render(k)
		aStyle := lipgloss.NewStyle().
			Foreground(m.currentTheme.Text).
			Render(" " + a)
		return kStyle + aStyle
	}

	row1 := fmt.Sprintf("%s  %s  %s  %s", badge("↑↓", "Nav"), badge("/", "Search"), badge("t", "Theme"), badge("s", "Snapshot"))
	row2 := fmt.Sprintf("%s  %s  %s  %s", badge("?", "Help"), badge("K", "Kill"), badge(":", "Command"), badge("q", "Quit"))
	if m.searchQuery != "" {
		row2 += "  " + badge("Esc", "Clear")
	}

	status := fmt.Sprintf("Theme: %s • Sort: %s • %d conns", m.currentTheme.ID, m.sortMode, len(m.allRows))
	if m.statusMessage != "" {
		status = "● " + m.statusMessage
	}
	statusLine := lipgloss.NewStyle().Foreground(m.currentTheme.Status).Render(status)

	body := title + "\n" + row1 + "\n" + row2 + "\n" + statusLine
	return cardStyle.Render(body)
}
