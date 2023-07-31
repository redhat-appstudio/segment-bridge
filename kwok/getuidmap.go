package kwok

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

const (
	contextName      = "kwok"
	clusterServerURL = "http://127.0.0.1:8080"
	scriptPath       = "../scripts/get-uid-map.sh"

	FailedToSetClusterMsg = "failed to set cluster: %w"
	FailedToSetContextMsg = "failed to set context: %w"
	FailedToUseContextMsg = "failed to use context: %w"
)

type ScriptExecutor interface {
	Execute(scriptPath string) ([]byte, error)
}

type ShellScriptExecutor struct{}

func (s *ShellScriptExecutor) Execute(scriptPath string) ([]byte, error) {
	return exec.Command("/bin/sh", scriptPath).CombinedOutput()
}

func runOcCommand(args ...string) error {
	cmd := exec.Command("oc", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run oc command: %w", err)
	}
	return nil
}

func setCluster(name, server string) error {
	return runOcCommand("config", "set-cluster", name, "--server="+server)
}

func setContext(name, cluster string) error {
	return runOcCommand("config", "set-context", name, "--cluster="+cluster)
}

func useContext(name string) error {
	return runOcCommand("config", "use-context", name)
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

func setUpClusterConfiguration() (map[string]int64, error) {
	actions := []func() error{
		func() error { return setCluster(contextName, clusterServerURL) },
		func() error { return setContext(contextName, contextName) },
		func() error { return useContext(contextName) },
	}

	for _, action := range actions {
		if err := action(); err != nil {
			return nil, err
		}
	}

	executor := &ShellScriptExecutor{}
	m, err := ExecuteAndParseScript(executor, scriptPath)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func validateMap(m map[string]int64) bool {
	if len(m) == 0 {
		return false
	}

	for user := range m {
		if user == "<no value>" || user == "" {
			return false
		}
	}
	return true
}
