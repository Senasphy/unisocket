package snapshot

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/table"
	"os"
	"encoding/json"
	"fmt"
	"time"
)




type portSnapshot struct{
	TimeStamp string `json:"timestamp"`
	Entries []portEntry `json:"entries"`
}

type portEntry struct{
	Port string `json:"port"`
	ProcessId string `json:"process-id"`
	ProcessName string `json:"process-name"`
	Protocol string `json:"protocol"`
	State string `json:"state"`
}

type SavedMessageString string

func SaveStateToFile(rows []table.Row, filename string) tea.Cmd{

	return func () tea.Msg{
		portsFile, err := os.Create(filename)
		var portEntries []portEntry
		if err != nil{
			return SavedMessageString("Can not create file")
		}
		defer portsFile.Close()

		for _, row := range rows{
			if len(row) < 5{
				continue
			}
			portEntries = append(portEntries, portEntry{
				Port: row[0],
				ProcessName: row[1],
				ProcessId: row[2],
				Protocol: row[3],
				State: row[4],
			})

		}
		snapshot := portSnapshot{
			TimeStamp: time.Now().Format(time.RFC3339),
			Entries: portEntries,
		};
		snapshotJson, err := json.MarshalIndent(snapshot, "", " ")
		if err != nil{
			return SavedMessageString("Error creating the Json file")
		}
		
		err = os.WriteFile(filename, snapshotJson, 0644)
		if err != nil{
			return SavedMessageString("Error saving the information to a file")
		}
		return SavedMessageString(fmt.Sprintf("JSON Successfully saved to file. You can find it at %s in the same directory.", filename))
	}


}
