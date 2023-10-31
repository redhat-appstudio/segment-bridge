package kwokfixture

import (
	"fmt"
	"os/exec"
)

const (
	contextName      = "kwok"
	clusterServerURL = "http://127.0.0.1:8080"
)

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

func SetUpClusterConfiguration() error {
	configurationActions := []func() error{
		func() error { return setCluster(contextName, clusterServerURL) },
		func() error { return setContext(contextName, contextName) },
		func() error { return useContext(contextName) },
	}

	for _, action := range configurationActions {
		if err := action(); err != nil {
			return err
		}
	}

	return nil
}
