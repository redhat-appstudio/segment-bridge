package scripts

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// A version of exec.LookPath that can find our scripts
// Current implementation works by manipulating $PATH to include the directory
// where this Go file is located, assuming it is placed in the same location as
// the scripts
func LookPath(file string) (string, error) {
	_, goFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to find path of scripts via Go file name")
	}
	if err := pushToPath(path.Dir(goFile)); err != nil {
		return "", err
	}
	return exec.LookPath(file)
}

// Push the given directory in front of $PATH unless its already listed there
func pushToPath(dir string) error {
	osPath := os.Getenv("PATH")
	osPathList := filepath.SplitList(osPath)
	for _, pathDir := range osPathList {
		if dir == pathDir {
			return nil
		}
	}
	newOsPath := fmt.Sprintf("%s%c%s", dir, filepath.ListSeparator, osPath)
	return os.Setenv("PATH", newOsPath)
}

func GetRepoRootDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to find the path of the root directory")
	}
	dirPath := filepath.Dir(filepath.Dir(filename))
	return dirPath, nil
}

func ExecuteScript(t *testing.T, scriptPath string) []byte {
	cmd := exec.Command(scriptPath)
	output, err := cmd.Output()
	assert.NoError(t, err, "failed to run script")
	return output
}
