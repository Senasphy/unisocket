package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Senasphy/unisocket/internal/app"
	"github.com/Senasphy/unisocket/internal/config"
	"github.com/Senasphy/unisocket/internal/connections"
	"github.com/Senasphy/unisocket/internal/theme"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	showTCP         = flag.Bool("tcp", false, "List TCP connections only")
	showUDP         = flag.Bool("udp", false, "Show UDP connections only")
	showEstablished = flag.Bool("estab", false, "Show established connections only")
	showListening   = flag.Bool("lsg", false, "Show listening connections only")
	findByName      = flag.String("find", "", "List all connections with the specified name")
	themeName       = flag.String("theme", "", "Set color theme (default, catppuccin, tokyonight, gruvbox, dracula, nord, onedark, rose-pine, solarized, monokai, cyberpunk)")
	configPath      = flag.String("config", "", "Path to configuration file (default: ~/.config/unisocket/config.json)")
	printConfig     = flag.Bool("print-config", false, "Print current configuration JSON and exit")
)

func main() {
	flag.Parse()

	cfg := config.Load(*configPath)

	if *printConfig {
		data, _ := json.MarshalIndent(cfg, "", "  ")
		fmt.Println(string(data))
		os.Exit(0)
	}

	// Track which CLI flags were explicitly provided
	passed := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		passed[f.Name] = true
	})

	if passed["theme"] {
		cfg.Theme = *themeName
	}
	if passed["tcp"] {
		cfg.Filters.TCP = *showTCP
	}
	if passed["udp"] {
		cfg.Filters.UDP = *showUDP
	}
	if passed["estab"] {
		cfg.Filters.Established = *showEstablished
	}
	if passed["lsg"] {
		cfg.Filters.Listening = *showListening
	}
	if passed["find"] {
		cfg.Filters.FindByName = *findByName
	}

	activeTheme := theme.Default()
	if th, ok := theme.Get(cfg.Theme); ok {
		activeTheme = th
	} else if cfg.Theme != "" {
		fmt.Printf("Warning: Unknown theme '%s', using default.\n", cfg.Theme)
	}

	filters := connections.Filters{
		TCP:         cfg.Filters.TCP,
		UDP:         cfg.Filters.UDP,
		Established: cfg.Filters.Established,
		Listening:   cfg.Filters.Listening,
		FindByName:  cfg.Filters.FindByName,
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
	t.SetStyles(activeTheme.TableStyles())

	m := app.New(t, filters, cfg, activeTheme)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
