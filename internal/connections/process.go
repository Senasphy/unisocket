package connections

import (
	"fmt"
	"os/exec"
	"runtime"
)

func KillProcess(pid string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/PID", pid, "/F")
	default:
		cmd = exec.Command("kill", "-TERM", pid)
	}

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return nil
}
