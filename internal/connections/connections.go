package connections

import (
	"fmt"
	"os/exec"
	"strings"
	"github.com/charmbracelet/bubbles/table"
	"os"
	"encoding/json"
	_ "embed"
)


type Filters struct{
	TCP bool
	UDP bool
	Established bool
	Listening bool
	FindByName string
}


//go:embed ports.json
var portsJSON []byte

var portMap = make(map[string]string)

func init() {
    if err := json.Unmarshal(portsJSON, &portMap); err != nil {
        fmt.Fprintf(os.Stderr, " Failed to parse embedded ports.json: %v\n", err)
    }
}

func GetConnections(filter Filters) []table.Row{
	var rows []table.Row
		ssString := []string{"-anp", "-w"}
		if filter.TCP {
			ssString = append(ssString, "-t")
		}
		if filter.UDP {
			ssString = append(ssString, "-u")
		}
		if !filter.TCP && !filter.UDP {
			ssString = append(ssString, "-tu")
		}

		ssCmd := exec.Command("sudo", append([]string{"ss"}, ssString...)...)
		out, err := ssCmd.Output()
		if err != nil {
			return nil
		}



			
		
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for i, line := range lines {
			if i == 0 {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}

			protocol := fields[0]
			state := fields[1]
			addrPort := fields[4]
			idx := strings.LastIndex(addrPort, ":")
			if idx == -1 { continue }
			port := addrPort[idx+1:]
			// Apply service naming from map
			if service, exists := portMap[port]; exists {
				port = fmt.Sprintf("%s (%s)", port, service)
			}

			if filter.Established && !strings.EqualFold(state, "ESTAB") {
				continue
			}
			if filter.Listening && !strings.EqualFold(state, "LISTEN") {
				continue
			}

			if len(fields) > 6 {
				processInfo := fields[6]
				processName, processId := getDetailedProcess(processInfo)

				if filter.FindByName != "" {
					nameMatch := strings.Contains(strings.ToLower(processName), strings.ToLower(filter.FindByName))
					serviceMatch := strings.Contains(strings.ToLower(port), strings.ToLower(filter.FindByName))
					if !nameMatch && !serviceMatch {
						continue
					}
				}
				rows = append(rows, table.Row{port, processId, processName, protocol, state})
			} else {
				if filter.FindByName != "" {
					continue
				}
				 rows = append(rows, table.Row{port, "-", "-", protocol, state}) 
			}
		}
		return rows


	}



func getDetailedProcess(raw string) (string, string) {
	if !strings.Contains(raw, "\"") || !strings.Contains(raw, "pid=") {
		return "Unknown", "-"
	}

	nameStart := strings.Index(raw, "\"") + 1
	relativeNameEnd := strings.Index(raw[nameStart:], "\"")
	processName := raw[nameStart : nameStart+relativeNameEnd]

	pIdStart := strings.Index(raw, "pid=") + 4
	relativePIdEnd := strings.Index(raw[pIdStart:], ",")

	if relativePIdEnd == -1 {
		relativePIdEnd = strings.Index(raw[pIdStart:], ")")
	}
	pId := raw[pIdStart : pIdStart+relativePIdEnd]
	return processName, pId
}


