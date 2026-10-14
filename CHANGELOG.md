# Changelog

All notable changes to this project will be documented in this file.

## [v1.0.2] - 2026-10-14

### Added
- Multi-theme engine with 11 built-in color schemes (Default, Catppuccin Mocha, Tokyo Night, Gruvbox Dark, Dracula, Nord, One Dark, Rosé Pine, Solarized Dark, Monokai Pro, Cyberpunk Neon).
- Interactive Theme Switcher accessible via `t` keybinding or `:theme` command, featuring live theme preview and persistence.
- JSON configuration support located at `~/.config/unisocket/config.json` with support for custom theme definitions, default sorting, snapshot filename, filters, and refresh rate.
- CLI flags: `-theme` to specify an initial theme, `-config` to supply a custom config file path, and `-print-config` to inspect effective settings.
- Help and Command Reference modal accessible via `?`, `h`, or `:help` / `:h`.
- Alternate screen buffer mode (`tea.WithAltScreen()`) to preserve terminal history upon exit.

### Changed
- Command row `:theme [name]` or `:t [name]` allows direct theme switching with alias support.
- Search input navigation ergonomics: `Esc` and `Ctrl+C` exit search mode without terminating the application; `Tab` and `Down Arrow` commit filter and transfer cursor focus to the connection table.

## [v1.0.1] - 2026-05-26

### Added
- One-command install scripts for Linux/macOS (`install.sh`) and Windows (`install.ps1`).
- Automated release workflows with GoReleaser.

### Changed
- Refactored internal architecture for process monitoring and socket inspection.
- Enhanced table view with live connection metrics.

## [v1.0.0] - 2026-05-26

### Added
- Initial cross-platform release supporting Linux, macOS, and Windows.
- Process-aware socket monitoring (PID, process name, local port, protocol, state).
- JSON snapshot export functionality (`s` key).
- Process termination mechanism (`K` key).
- In-table search and filtering (`/` key).
- Column sorting by port (`:sn`), process name (`:spn`), and PID (`:spid`).
