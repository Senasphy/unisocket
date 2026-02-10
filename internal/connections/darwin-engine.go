//go:build darwin

package connections

import (
	"os/exec"
	"strings"
	"github.com/charmbracelet/bubbles/table"
)

func GetConnections(filter Filters) []table.Row {
	var rows []table.Row
	out, err := exec.Command("sudo", "lsof", "-i", "-n", "-P").Output()
	if err != nil { return nil }

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i, line := range lines {
		if i == 0 { continue }
		f := strings.Fields(line)
		if len(f) < 9 { continue }

		proto := strings.ToLower(f[7])
		if (filter.TCP && !strings.Contains(proto, "tcp")) || (filter.UDP && !strings.Contains(proto, "udp")) { continue }

		nameField := f[8]
		idx := strings.LastIndex(nameField, ":")
		if idx == -1 { continue }
		
		portPart := nameField[idx+1:]
		state := "N/A"
		if sIdx := strings.Index(portPart, "("); sIdx != -1 {
			state = portPart[sIdx+1 : strings.Index(portPart, ")")]
			portPart = portPart[:sIdx]
		}
		port := formatPort(portPart)

		if (filter.Established && !strings.EqualFold(state, "ESTABLISHED")) ||
		   (filter.Listening && !strings.EqualFold(state, "LISTEN")) {
			continue
		}

		if filter.FindByName != "" && !strings.Contains(strings.ToLower(f[0]+port), strings.ToLower(filter.FindByName)) {
			continue
		}
		rows = append(rows, table.Row{port, f[1], f[0], proto, state})
	}
	return rows
}
