//go:build !windows

package steamcmd

import (
	"os/exec"
)

func setSysProcAttr(cmd *exec.Cmd) {
	// No-op for non-Windows platforms
}
