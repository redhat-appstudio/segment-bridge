package kwok

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed kwok_container_template.tmpl
var KwokServiceManifest string

// GetKubeconfig get the full path to the kwok directory's
// kubeconfig file, no matter where the functions is called from,
// then sets it as the "KUBECONFIG" environment variable.
func SetKubeconfig() error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to find the path of the current script")
	}
	dirPath := filepath.Dir(filename)
	kubeconfigPath := filepath.Join(dirPath, "kubeconfig")
	os.Setenv("KUBECONFIG", kubeconfigPath)

	return nil
}
