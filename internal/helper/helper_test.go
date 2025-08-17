package helper

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestStrToInt64(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		panicIfErr bool
		want       int64
		wantPanic  bool
	}{
		{"empty string", "", false, 0, false},
		{"zero", "0", false, 0, false},
		{"positive number", "123", false, 123, false},
		{"negative number", "-123", false, -123, false},
		{"with whitespace", " 456 ", false, 456, false},
		{"invalid input no panic", "abc", false, 0, false},
		{"invalid input with panic", "abc", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("StrToInt64(%q, %v) should have panicked", tt.input, tt.panicIfErr)
					}
				}()
			}

			got := StrToInt64(tt.input, tt.panicIfErr)
			if !tt.wantPanic && got != tt.want {
				t.Errorf("StrToInt64(%q, %v) = %d; want %d", tt.input, tt.panicIfErr, got, tt.want)
			}
		})
	}
}

func TestStrToInt32(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		panicIfErr bool
		want       int32
		wantPanic  bool
	}{
		{"empty string", "", false, 0, false},
		{"zero", "0", false, 0, false},
		{"positive number", "123", false, 123, false},
		{"negative number", "-123", false, -123, false},
		{"with whitespace", " 456 ", false, 456, false},
		{"invalid input no panic", "abc", false, 0, false},
		{"invalid input with panic", "abc", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("StrToInt32(%q, %v) should have panicked", tt.input, tt.panicIfErr)
					}
				}()
			}

			got := StrToInt32(tt.input, tt.panicIfErr)
			if !tt.wantPanic && got != tt.want {
				t.Errorf("StrToInt32(%q, %v) = %d; want %d", tt.input, tt.panicIfErr, got, tt.want)
			}
		})
	}
}

func TestConvertToFloat(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		panicIfErr bool
		want       float64
		wantPanic  bool
	}{
		{"empty string", "", false, 0, false},
		{"zero", "0", false, 0, false},
		{"positive number", "123.45", false, 123.45, false},
		{"negative number", "-123.45", false, -123.45, false},
		{"with whitespace", " 456.78 ", false, 456.78, false},
		{"integer as float", "789", false, 789, false},
		{"invalid input no panic", "abc", false, 0, false},
		{"invalid input with panic", "abc", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("ConvertToFloat(%q, %v) should have panicked", tt.input, tt.panicIfErr)
					}
				}()
			}

			got := ConvertToFloat(tt.input, tt.panicIfErr)
			if !tt.wantPanic && got != tt.want {
				t.Errorf("ConvertToFloat(%q, %v) = %f; want %f", tt.input, tt.panicIfErr, got, tt.want)
			}
		})
	}
}

func TestGetFileSize(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want string
	}{
		{"zero bytes", 0, "0.00 Kb"},
		{"bytes", 100, "0.10 Kb"},
		{"kilobytes", 1024, "1.00 Kb"},
		{"kilobytes decimal", 1536, "1.50 Kb"},
		{"megabytes", 1048576, "1.00 Mb"}, // 1024*1024
		{"megabytes decimal", 1572864, "1.50 Mb"}, // 1.5*1024*1024
		{"gigabytes", 1073741824, "1.00 Gb"}, // 1024*1024*1024
		{"gigabytes decimal", 1610612736, "1.50 Gb"}, // 1.5*1024*1024*1024
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFileSize(tt.size)
			if got != tt.want {
				t.Errorf("GetFileSize(%d) = %q; want %q", tt.size, got, tt.want)
			}
		})
	}
}

func TestMemoryUsage(t *testing.T) {
	// This is a basic test to ensure the function returns a string
	// with the expected format. We can't test exact values as they'll vary.
	result := MemoryUsage()

	// Check that the result contains the expected substrings
	if !strings.Contains(result, "TotalAlloc:") || !strings.Contains(result, "MB") || !strings.Contains(result, "Sys:") {
		t.Errorf("MemoryUsage() = %q; expected to contain 'TotalAlloc:', 'Sys:', and 'MB'", result)
	}
}

func TestRuntimeStatistics(t *testing.T) {
	// Create a temporary file for testing
	tempFile := "temp_test_file.txt"
	content := []byte("test content")
	err := os.WriteFile(tempFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile) // Clean up after test

	startTime := time.Now().Add(-time.Second) // Simulate that the function started 1 second ago
	result := RuntimeStatistics(startTime, tempFile)

	// Check that the result contains expected parts
	if !strings.Contains(result, tempFile) ||
		!strings.Contains(result, "Processed") ||
		!strings.Contains(result, "TotalAlloc:") {
		t.Errorf("RuntimeStatistics() = %q; expected to contain file name, 'Processed', and memory stats", result)
	}
}

func TestAppHelp(t *testing.T) {
	// We can only test the false case safely without mocking os.Exit
	// For the true case, we would need to mock os.Exit which is not possible
	// in a standard way in Go without modifying the source code.

	// Capture log output
	var buf bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(oldOutput)

	// Test with help = false (should not exit or log)
	AppHelp(false)

	// Verify no output was logged
	if buf.Len() > 0 {
		t.Errorf("AppHelp(false) wrote to log, which it shouldn't: %s", buf.String())
	}

	// Note: We cannot safely test the true case (AppHelp(true))
	// as it calls os.Exit directly, which would terminate the test process.
	// In a real-world scenario, we would refactor the code to make it more testable
	// by injecting the exit function or returning a value instead of calling os.Exit directly.
}

// Benchmarks
func BenchmarkStrToInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrToInt64("12345", false)
	}
}

func BenchmarkStrToInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrToInt32("12345", false)
	}
}

func BenchmarkConvertToFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConvertToFloat("123.45", false)
	}
}

func BenchmarkGetFileSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFileSize(1048576) // 1MB
	}
}
