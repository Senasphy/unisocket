//go:build windows

package connections

import (
	"os/exec"
	"strings"
	"github.com/charmbracelet/bubbles/table"
)

func GetConnections(filter Filters) []table.Row {
	var rows []table.Row
	// -a (all), -n (numeric), -o (shows PID)
	out, err := exec.Command("netstat", "-ano").Output()
	if err != nil { return nil }

	lines := strings.Split(string(out), "\r\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		// Windows netstat output starts with "TCP" or "UDP"
		if len(fields) < 4 { continue }
		proto := fields[0]
		if proto != "TCP" && proto != "UDP" { continue }

		// Filter protocol
		if filter.TCP && proto != "TCP" { continue }
		if filter.UDP && proto != "UDP" { continue }

		localAddr := fields[1]
		idx := strings.LastIndex(localAddr, ":")
		if idx == -1 { continue }
		port := formatPort(localAddr[idx+1:])

		state := ""
		pid := ""
		if proto == "TCP" {
			state = fields[3]
			pid = fields[4]
		} else {
			// UDP is stateless in netstat, so PID is in index 3
			state = "UDP"
			pid = fields[3]
		}

		if filter.Established && !strings.EqualFold(state, "ESTABLISHED") { continue }
		if filter.Listening && !strings.EqualFold(state, "LISTENING") { continue }

		// Note: Windows netstat doesn't give process names, only PIDs.
		// To get names, you'd need another call to 'tasklist', but for now:
		processName := "PID: " + pid 

		if filter.FindByName != "" {
			if !strings.Contains(strings.ToLower(port), strings.ToLower(filter.FindByName)) {
				continue
			}
		}
		rows = append(rows, table.Row{port, pid, processName, proto, state})
	}
	return rows
}
