package testfixture

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type ScriptExecutor interface {
	Execute(scriptPath string) ([]byte, error)
}

func ExecuteAndParseScript(executor ScriptExecutor, scriptPath string) (map[string]int64, error) {
	m := make(map[string]int64)

	output, err := executor.Execute(scriptPath)
	if err != nil {
		return m, fmt.Errorf("failed to execute command\nOutput: %s\nError: %w", string(output), err)
	}

	if err := json.Unmarshal(output, &m); err != nil {
		return m, fmt.Errorf("failed to parse result: %w", err)
	}

	return m, nil
}

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
