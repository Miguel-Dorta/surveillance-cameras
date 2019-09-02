package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Executes "mv" command creating the necessary directories in destiny
func Move(origin, destiny string) error {
	if err := os.MkdirAll(filepath.Dir(destiny), 0755); err != nil {
		return fmt.Errorf("cannot create parent directories for path \"%s\": %s", destiny, err)
	}

	cmd := exec.Command("mv", origin, destiny)
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		return fmt.Errorf("error moving file from \"%s\" to \"%s\": %s. Output: %s", origin, destiny, err, string(out))
	}
	return nil
}
