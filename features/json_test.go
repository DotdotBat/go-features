package features

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Test function
func TestReadConfigFromFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		fileContent string
		fileName    string
		want        *JSONConfig
		wantErr     bool
	}{
		{
			name:     "Valid JSON file",
			fileName: "valid_config.json",
			fileContent: `{
				"name": "TestApp",
				"version": "1.0.0",
				"port": 8080,
				"enabled": true,
				"features": ["auth", "logging", "metrics"]
			}`,
			want: &JSONConfig{
				Name:     "TestApp",
				Version:  "1.0.0",
				Port:     8080,
				Enabled:  true,
				Features: []string{"auth", "logging", "metrics"},
			},
			wantErr: false,
		},
		{
			name:        "Invalid JSON format",
			fileName:    "invalid_config.json",
			fileContent: `{"name": "TestApp", "version": }`, // Invalid JSON
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "Empty JSON file",
			fileName:    "empty_config.json",
			fileContent: `{}`,
			want: &JSONConfig{
				Name:     "",
				Version:  "",
				Port:     0,
				Enabled:  false,
				Features: nil,
			},
			wantErr: false,
		},
		{
			name:        "File does not exist",
			fileName:    "nonexistent.json",
			fileContent: "", // Won't create this file
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file (except for "file does not exist" test)
			if tt.name != "File does not exist" {
				filePath := filepath.Join(tempDir, tt.fileName)
				err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			// Test the function
			filePath := filepath.Join(tempDir, tt.fileName)
			got, err := ReadConfigFromFile(filePath)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check result
			if !tt.wantErr && !compareConfigs(got, tt.want) {
				t.Errorf("ReadConfigFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare Config structs
func compareConfigs(a, b *JSONConfig) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Name != b.Name || a.Version != b.Version || a.Port != b.Port || a.Enabled != b.Enabled {
		return false
	}
	if len(a.Features) != len(b.Features) {
		return false
	}
	for i := range a.Features {
		if a.Features[i] != b.Features[i] {
			return false
		}
	}
	return true
}

// Additional test for concurrent reads
func TestReadConfigFromFileConcurrent(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "concurrent_test.json")

	// Create a test file
	content := `{"name": "ConcurrentTest", "version": "1.0.0", "port": 9090, "enabled": true, "features": ["test"]}`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Run multiple goroutines reading the same file
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			config, err := ReadConfigFromFile(testFile)
			if err != nil {
				t.Errorf("Concurrent read failed: %v", err)
			}
			if config.Name != "ConcurrentTest" {
				t.Errorf("Expected name 'ConcurrentTest', got '%s'", config.Name)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Benchmark test
func BenchmarkReadConfigFromFile(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "benchmark_test.json")

	content := `{"name": "BenchmarkTest", "version": "1.0.0", "port": 8080, "enabled": true, "features": ["auth", "logging", "metrics"]}`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ReadConfigFromFile(testFile)
		if err != nil {
			b.Fatalf("Benchmark read failed: %v", err)
		}
	}
}

// Helper to load config from disk for test verification
func loadConfigFromFile(t *testing.T, path string) JSONConfig {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	var cfg JSONConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("could not unmarshal JSON: %v", err)
	}
	return cfg
}

func TestWriteConfig_CreatesDefaultIfNoneExists(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	WriteToConfig(configPath, "Name", "DefaultApp")

	cfg := loadConfigFromFile(t, configPath)

	if cfg.Name != "DefaultApp" {
		t.Errorf("expected Name to be 'DefaultApp', got %s", cfg.Name)
	}
}

func TestWriteConfig_UpdatesSingleFieldWithoutOverwritingOthers(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Start with an initial config file
	initialConfig := JSONConfig{
		Name:     "MyApp",
		Version:  "1.0",
		Port:     8080,
		Enabled:  true,
		Features: []string{"feature1", "feature2"},
	}
	data, _ := json.MarshalIndent(initialConfig, "", "  ")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	WriteToConfig(configPath, "Port", 9090)

	cfg := loadConfigFromFile(t, configPath)

	// Check updated field
	if cfg.Port != 9090 {
		t.Errorf("expected Port 9090, got %d", cfg.Port)
	}

	// Ensure other fields remain unchanged
	expected := initialConfig
	expected.Port = 9090
	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("config mismatch:\nexpected %+v\n got %+v", expected, cfg)
	}
}
