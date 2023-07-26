package scripts

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
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
