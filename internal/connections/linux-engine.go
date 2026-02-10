//go:build linux

package connections

import (
	"os/exec"
	"strings"
	"github.com/charmbracelet/bubbles/table"
)

func GetConnections(filter Filters) []table.Row {
	var rows []table.Row
	args := []string{"-anp", "-H"} // -H skips the header row

	if filter.TCP { args = append(args, "-t") }
	if filter.UDP { args = append(args, "-u") }
	if !filter.TCP && !filter.UDP { args = append(args, "-tu") }

	out, err := exec.Command("sudo", append([]string{"ss"}, args...)...).Output()
	if err != nil { return nil }

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 { continue }

		protocol := fields[0]
		state := fields[1]
		addrPort := fields[4]
		
		idx := strings.LastIndex(addrPort, ":")
		if idx == -1 { continue }
		port := formatPort(addrPort[idx+1:])

		if filter.Established && !strings.EqualFold(state, "ESTAB") { continue }
		if filter.Listening && !strings.EqualFold(state, "LISTEN") { continue }

		processName, processId := "Unknown", "-"
		if len(fields) > 6 {
			processName, processId = parseLinuxProcess(fields[6])
		}

		if filter.FindByName != "" {
			if !strings.Contains(strings.ToLower(processName), strings.ToLower(filter.FindByName)) &&
			   !strings.Contains(strings.ToLower(port), strings.ToLower(filter.FindByName)) {
				continue
			}
		}
		rows = append(rows, table.Row{port, processId, processName, protocol, state})
	}
	return rows
}

func parseLinuxProcess(raw string) (string, string) {
	if !strings.Contains(raw, "\"") || !strings.Contains(raw, "pid=") {
		return "Unknown", "-"
	}
	nameStart := strings.Index(raw, "\"") + 1
	nameEnd := strings.Index(raw[nameStart:], "\"")
	pIdStart := strings.Index(raw, "pid=") + 4
	pIdEnd := strings.Index(raw[pIdStart:], ",")
	if pIdEnd == -1 { pIdEnd = strings.Index(raw[pIdStart:], ")") }
	
	return raw[nameStart : nameStart+nameEnd], raw[pIdStart : pIdStart+pIdEnd]
}
