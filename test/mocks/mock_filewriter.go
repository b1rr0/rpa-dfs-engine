package mocks

import (
	"fmt"
	"os"
)

type MockFileWriter struct {
	WrittenFiles    map[string][]byte
	ShouldFail      bool
	ErrorMessage    string
	FailOnFilenames []string // Specific filenames that should cause failure
}

func NewMockFileWriter() *MockFileWriter {
	return &MockFileWriter{
		WrittenFiles: make(map[string][]byte),
	}
}

func (m *MockFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if m.ShouldFail {
		if m.ErrorMessage != "" {
			return fmt.Errorf(m.ErrorMessage)
		}
		return fmt.Errorf("mock write failure")
	}

	for _, failFile := range m.FailOnFilenames {
		if filename == failFile {
			return fmt.Errorf("permission denied for file: %s", filename)
		}
	}

	m.WrittenFiles[filename] = data
	return nil
}

func (m *MockFileWriter) GetWrittenContent(filename string) []byte {
	return m.WrittenFiles[filename]
}

func (m *MockFileWriter) WasFileWritten(filename string) bool {
	_, exists := m.WrittenFiles[filename]
	return exists
}

func (m *MockFileWriter) GetWrittenFileCount() int {
	return len(m.WrittenFiles)
}

func (m *MockFileWriter) Reset() {
	m.WrittenFiles = make(map[string][]byte)
	m.ShouldFail = false
	m.ErrorMessage = ""
	m.FailOnFilenames = nil
}
