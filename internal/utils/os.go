package utils

import (
	"os/exec"
	"runtime"
)

// OpenDirectory opens the specified directory in the file explorer based on the operating system
func OpenDirectory(path string) error {
	var cmd *exec.Cmd
	switch os := runtime.GOOS; os {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
