package testfixture

import (
	"fmt"
	"os"
	"os/exec"
)

func RunScriptWithInputFile(filePath, scriptPath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	cmd := exec.Command(scriptPath)
	cmd.Stdin = file

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing script: %w", err)
	}

	return output, nil
}
