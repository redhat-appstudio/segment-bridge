package testfixture

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Execute(scriptPath string) ([]byte, error) {
	args := m.Called(scriptPath)
	return args.Get(0).([]byte), args.Error(1)
}

func createTempFile(t *testing.T, content string) string {
	t.Helper()
	file, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	return file.Name()
}

func createTempScript(t *testing.T, content string) string {
	t.Helper()
	scriptFile, err := os.CreateTemp("", "*.sh")
	if err != nil {
		t.Fatalf("Failed to create temp script: %v", err)
	}
	defer scriptFile.Close()

	if _, err := scriptFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp script: %v", err)
	}

	if err := os.Chmod(scriptFile.Name(), 0700); err != nil {
		t.Fatalf("Failed to make script executable: %v", err)
	}

	return scriptFile.Name()
}

func TestRunScriptWithInputFile(t *testing.T) {
	t.Run("SuccessfulExecution", func(t *testing.T) {
		inputFilePath := createTempFile(t, "Hello world")
		defer os.Remove(inputFilePath)

		scriptFilePath := createTempScript(t, "#!/bin/sh\ncat")
		defer os.Remove(scriptFilePath)

		output, err := RunScriptWithInputFile(inputFilePath, scriptFilePath)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expectedOutput := "Hello world"
		if string(output) != expectedOutput {
			t.Errorf("Expected %q, got %q", expectedOutput, output)
		}
	})

	t.Run("ScriptExecutionError", func(t *testing.T) {
		inputFilePath := createTempFile(t, "")
		defer os.Remove(inputFilePath)

		scriptFilePath := createTempScript(t, "#!/bin/sh\nexit 1")
		defer os.Remove(scriptFilePath)

		if _, err := RunScriptWithInputFile(inputFilePath, scriptFilePath); err == nil {
			t.Errorf("Expected an error, got none")
		}
	})

	t.Run("FileOpenError", func(t *testing.T) {
		scriptFilePath := createTempScript(t, "#!/bin/sh\ncat")
		defer os.Remove(scriptFilePath)

		if _, err := RunScriptWithInputFile("/tmp/ewefewg34234", scriptFilePath); err == nil {
			t.Errorf("Expected an error, got none")
		}
	})
}

func TestExecuteAndParseScript(t *testing.T) {
	executor := new(MockExecutor)

	t.Run("Successful Execution", func(t *testing.T) {
		executor.On("Execute", "valid_script.sh").Return([]byte(`{"key": 123}`), nil)
		result, err := ExecuteAndParseScript(executor, "valid_script.sh")
		assert.NoError(t, err)
		assert.Equal(t, map[string]int64{"key": 123}, result)
		executor.AssertExpectations(t)
	})

	t.Run("Execution Failure", func(t *testing.T) {
		executor.On("Execute", "invalid_script.sh").Return([]byte{}, fmt.Errorf("execution failed"))
		_, err := ExecuteAndParseScript(executor, "invalid_script.sh")
		assert.Error(t, err)
		executor.AssertExpectations(t)
	})

	t.Run("JSON Parsing Failure", func(t *testing.T) {
		executor.On("Execute", "invalid_json.sh").Return([]byte("invalid json"), nil)
		_, err := ExecuteAndParseScript(executor, "invalid_json.sh")
		assert.Error(t, err)
		executor.AssertExpectations(t)
	})
}
