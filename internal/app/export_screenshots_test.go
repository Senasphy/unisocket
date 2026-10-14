package app

import (
	"os"
	"testing"

	"github.com/Senasphy/unisocket/internal/config"
	"github.com/Senasphy/unisocket/internal/connections"
	"github.com/Senasphy/unisocket/internal/theme"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func sampleRows() []table.Row {
	return []table.Row{
		{"53 (dns)", "1122", "systemd-resolve", "udp", "UNCONN"},
		{"53 (dns)", "1122", "systemd-resolve", "udp", "UNCONN"},
		{"68 (dhcp-client)", "1779", "NetworkManager", "udp", "ESTAB"},
		{"52105", "1660", "avahi-daemon", "udp", "UNCONN"},
		{"5353 (mdns)", "5505", "chrome", "udp", "UNCONN"},
		{"5353 (mdns)", "1660", "avahi-daemon", "udp", "UNCONN"},
		{"5353 (mdns)", "1660", "avahi-daemon", "udp", "UNCONN"},
		{"46319", "1660", "avahi-daemon", "udp", "UNCONN"},
		{"53 (dns)", "1122", "systemd-resolve", "tcp", "LISTEN"},
		{"631 (ipp)", "1958", "cupsd", "tcp", "LISTEN"},
		{"5432 (postgresql)", "2039", "postgres", "tcp", "LISTEN"},
		{"33034", "3624", "zen", "tcp", "ESTAB"},
		{"60814", "7718", "telegram-desktop", "tcp", "ESTAB"},
		{"46800", "3624", "zen", "tcp", "ESTAB"},
		{"53284", "3624", "zen", "tcp", "ESTAB"},
		{"53930", "3624", "zen", "tcp", "ESTAB"},
		{"51146", "-", "Unknown", "tcp", "FIN-WAIT-1"},
		{"53238", "5756", "chrome", "tcp", "ESTAB"},
		{"48270", "-", "Unknown", "tcp", "TIME-WAIT"},
		{"54290", "-", "Unknown", "tcp", "TIME-WAIT"},
		{"8080 (http-alt)", "14302", "node", "tcp", "LISTEN"},
		{"6379 (redis)", "2104", "redis-server", "tcp", "LISTEN"},
		{"27017 (mongodb)", "2240", "mongod", "tcp", "LISTEN"},
	}
}

func TestExportScreenshots(t *testing.T) {
	if os.Getenv("EXPORT_SCREENSHOTS") != "1" {
		t.Skip("skipping screenshot export; run with EXPORT_SCREENSHOTS=1")
	}

	lipgloss.SetColorProfile(termenv.TrueColor)

	w := 154
	h := 36

	rows := sampleRows()

	// 1. Tokyo Night theme
	{
		th, _ := theme.Get("tokyonight")
		cfg := config.DefaultConfig()
		cfg.Theme = "tokyonight"
		m := New(table.New(), connections.Filters{}, cfg, th).(Model)
		m.width = w
		m.height = h
		m.allRows = rows
		m.updateTableDimensions()
		m.applyRows()
		m.table.SetCursor(0)

		out := m.View()
		_ = os.WriteFile("/tmp/screen_tokyonight.ansi", []byte(out), 0644)
	}

	// 2. Catppuccin Mocha with side-by-side Theme Switcher
	{
		th, _ := theme.Get("catppuccin")
		cfg := config.DefaultConfig()
		cfg.Theme = "catppuccin"
		m := New(table.New(), connections.Filters{}, cfg, th).(Model)
		m.width = w
		m.height = h
		m.allRows = rows
		m.openThemePicker()
		// Cursor on catppuccin or tokyonight
		themes := theme.List()
		for i, item := range themes {
			if item.ID == "catppuccin" {
				m.themePickerIndex = i
				break
			}
		}
		m.updateTableDimensions()
		m.applyRows()
		m.table.SetCursor(1)

		out := m.View()
		_ = os.WriteFile("/tmp/screen_catppuccin_switcher.ansi", []byte(out), 0644)
	}

	// 3. Dracula with Help modal
	{
		th, _ := theme.Get("dracula")
		cfg := config.DefaultConfig()
		cfg.Theme = "dracula"
		m := New(table.New(), connections.Filters{}, cfg, th).(Model)
		m.width = w
		m.height = h
		m.allRows = rows
		m.isHelp = true
		m.updateTableDimensions()
		m.applyRows()
		m.table.SetCursor(2)

		out := m.View()
		_ = os.WriteFile("/tmp/screen_dracula_help.ansi", []byte(out), 0644)
	}

	// 4. Dracula standard table
	{
		th, _ := theme.Get("dracula")
		cfg := config.DefaultConfig()
		cfg.Theme = "dracula"
		m := New(table.New(), connections.Filters{}, cfg, th).(Model)
		m.width = w
		m.height = h
		m.allRows = rows
		m.updateTableDimensions()
		m.applyRows()
		m.table.SetCursor(3)

		out := m.View()
		_ = os.WriteFile("/tmp/screen_dracula.ansi", []byte(out), 0644)
	}

	// 5. Gruvbox with search filter
	{
		th, _ := theme.Get("gruvbox")
		cfg := config.DefaultConfig()
		cfg.Theme = "gruvbox"
		m := New(table.New(), connections.Filters{}, cfg, th).(Model)
		m.width = w
		m.height = h
		m.allRows = rows
		m.searchQuery = "chrome"
		m.isSearching = true
		m.searchInput.SetValue("chrome")
		m.updateTableDimensions()
		m.applyRows()

		out := m.View()
		_ = os.WriteFile("/tmp/screen_gruvbox.ansi", []byte(out), 0644)
	}
}
