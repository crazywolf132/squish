package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func RunTSC(srcDir, outDir, tsconfigPath string) error {
	args := []string{
		"tsc",
		"--declaration",
		"--emitDeclarationOnly",
		"--outDir", outDir,
	}

	if tsconfigPath != "" {
		args = append(args, "--project", tsconfigPath)
	} else {
		args = append(args, srcDir)
	}

	cmd := exec.Command("npx", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	Log("Running tsc command:", cmd.String())
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("tsc command failed: %w", err)
	}

	return nil
}
