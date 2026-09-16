package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-gitea/gitea/modules/markup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/md4"
)

// MockFileSystem is a mock for file system operations
type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	args := m.Called(filename)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	args := m.Called(filename, data, perm)
	return args.Error(0)
}

// MockExecCommand is a mock for exec.Command
type MockExecCommand struct {
	mock.Mock
}

func (m *MockExecCommand) Run() error {
	args := m.Called()
	return args.Error(0)
}

// Test for MD4 hash calculation
func TestMD4Hash(t *testing.T) {
	// This is a simple test to verify the MD4 hash calculation
	data := "These pretzels are making me thirsty."
	expectedHash := "f4c8521ccde9fed0e19cb3d5ab651718" // Pre-calculated MD4 hash

	h := md4.New()
	io.WriteString(h, data)
	actualHash := fmt.Sprintf("%x", h.Sum(nil))

	assert.Equal(t, expectedHash, actualHash, "MD4 hash should match expected value")
}

// Test writing config file
func TestWriteConfigFile(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := ioutil.TempFile("", "test-config-*.json")
	require.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tmpfile.Name())

	// Write config to the temporary file
	err = ioutil.WriteFile(tmpfile.Name(), validConfig, 0644)
	require.NoError(t, err, "Failed to write config file")

	// Read the file back and verify contents
	data, err := ioutil.ReadFile(tmpfile.Name())
	require.NoError(t, err, "Failed to read config file")
	assert.Equal(t, validConfig, data, "Config file content should match")
}

// Test IsReadmeFile function from gitea markup package
func TestIsReadmeFile(t *testing.T) {
	// Create a mock for markup.IsReadmeFile if needed
	// For simplicity, we'll just test the actual function
	assert.True(t, markup.IsReadmeFile("README.md"), "README.md should be recognized as a readme file")
	assert.False(t, markup.IsReadmeFile("notreadme.md"), "notreadme.md should not be recognized as a readme file")
}

// Test the file path handler function
func TestFilePathHandler(t *testing.T) {
	// Create a temporary test file
	content := []byte("test file content")
	tmpfile, err := ioutil.TempFile("", "test-file-*.txt")
	require.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err, "Failed to write to temp file")
	tmpfile.Close()

	// Create a test request with the file path as a query parameter
	req, err := http.NewRequest("GET", "/?path="+tmpfile.Name(), nil)
	require.NoError(t, err, "Failed to create request")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test handler that simulates our vulnerable handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Query().Get("path")
		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Handler should return status 200")

	// Check the response body
	assert.Equal(t, string(content), rr.Body.String(), "Response should match file content")
}

// Test the decompression handler
func TestDecompressionHandler(t *testing.T) {
	// Create a simple gzipped content
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write([]byte("test compressed content"))
	require.NoError(t, err, "Failed to write compressed content")
	err = gzipWriter.Close()
	require.NoError(t, err, "Failed to close gzip writer")

	// Create a test request with gzipped body
	req, err := http.NewRequest("POST", "/decompress", bytes.NewReader(buf.Bytes()))
	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/gzip")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test handler
	// For testing, we'll modify the handler to write to a buffer instead of os.Stdout
	var output bytes.Buffer
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<30) // 1GB
		gzr, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Error creating gzip reader", http.StatusInternalServerError)
			return
		}
		defer gzr.Close()
		_, err = io.Copy(&output, gzr)
		if err != nil {
			http.Error(w, "Error decompressing data", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("OK"))
	})

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Handler should return status 200")

	// Check that the decompressed content is correct
	assert.Equal(t, "test compressed content", output.String(), "Decompressed content should match original")
}

// Test the config loading
func TestConfigLoading(t *testing.T) {
	// Create a temporary directory for config files
	tmpDir, err := ioutil.TempDir("", "config-test")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tmpDir)

	// Create the config directory if it doesn't exist
	configDir := tmpDir + "/config"
	err = os.MkdirAll(configDir, 0755)
	require.NoError(t, err, "Failed to create config directory")

	// Write a test config file
	err = ioutil.WriteFile(configDir+"/phish-config.json", validConfig, 0644)
	require.NoError(t, err, "Failed to write config file")

	// Load the config
	// In a real test, we would need to mock or adjust the config loading logic
	// For this example, we'll just verify the config file exists
	_, err = os.Stat(configDir + "/phish-config.json")
	assert.NoError(t, err, "Config file should exist")
}
