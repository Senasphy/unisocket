package connections

import (
	"encoding/json"
	"fmt"
	"os"
	_ "embed"
)

type Filters struct {
	TCP         bool
	UDP         bool
	Established bool
	Listening   bool
	FindByName  string
}

//go:embed ports.json
var portsJSON []byte
var portMap = make(map[string]string)

func init() {
	if len(portsJSON) > 0 {
		if err := json.Unmarshal(portsJSON, &portMap); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse embedded ports.json: %v\n", err)
		}
	}
}

func formatPort(port string) string {
	if service, exists := portMap[port]; exists {
		return fmt.Sprintf("%s (%s)", port, service)
	}
	return port
}
