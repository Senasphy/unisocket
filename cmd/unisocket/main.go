package main

import (
	"flag"
	"fmt"
	"github.com/Senasphy/unisocket/internal/app"
	"github.com/Senasphy/unisocket/internal/connections"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var (
	showTCP         = flag.Bool("tcp", false, "List TCP connections only")
	showUDP         = flag.Bool("udp", false, "Show UDP connections only")
	showEstablished = flag.Bool("estab", false, "Show established connections only")
	showListening   = flag.Bool("lsg", false, "Show listening connections only")
	findByName      = flag.String("find", "", "List all connections with the specified name")
)

func main() {
	flag.Parse()

	filters := connections.Filters{
		TCP:         *showTCP,
		UDP:         *showUDP,
		Established: *showEstablished,
		Listening:   *showListening,
		FindByName:  *findByName,
	}
	columns := []table.Column{
		{Title: "PORT", Width: 22},
		{Title: "PROCESS ID", Width: 12},
		{Title: "PROCESS NAME", Width: 24},
		{Title: "PROTOCOL", Width: 10},
		{Title: "STATE", Width: 12},
	}

	// the rows should be populated with the result of the printTable
	rows := connections.GetConnections(filters)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Style the table header
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		BorderBottom(false).
		Padding(0, 1).
		MarginBottom(1).
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("24")).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("31")).
		Bold(true)
	t.SetStyles(s)

	m := app.New(t, filters)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
