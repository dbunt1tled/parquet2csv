# 🧠 AI Assistant Prompt for GoLang Project

You are an AI assistant working on a GoLang web service. Follow the rules below to ensure your code matches the project’s style and standards.

⸻

## 📐 Coding Style Guidelines
-	Follow idiomatic Go, inspired by Effective Go
-	Use camelCase for variables and PascalCase for exported functions/types
-	Keep functions small and focused
-	Structure: handlers → services → repositories
-	Error handling:
-	Wrap errors using ``errors.Wrap(err,"new context")``
-	Define and use custom error types when needed
-	Logging:
-	Prefer ``log/slog`` library, avoid third-party packages unless essential
-	Format code with gofmt standards (tabs, spacing)

⸻

## 🧪 Testing Guidelines
-	Use table-driven tests
-	Use Go’s built-in testing package
-	Avoid third-party testing frameworks
-	Organize tests in _test.go files alongside source files
-	Use subtests (t.Run(...)) for better test reporting
-	Cover edge cases in tests

Example of a unit test:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive numbers", 1, 2, 3},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

⸻

## 🚀 Benchmark Guidelines
-	Benchmark functions must start with Benchmark
-	Always use for i := 0; i < b.N; i++ loop for repeated execution
-	Move global setup before b.ResetTimer()
-	Use b.StopTimer() and b.StartTimer() for heavy setup/cleanup inside loop if needed
-	Avoid allocations and external influences during measurements
-	Cleanup after benchmark if necessary

Simple benchmark example:
```go
func BenchmarkAdd(b *testing.B) {
    a, c := 1, 2

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
    _ = Add(a, c)
    }
}
    
```
Advanced benchmark with per-iteration setup/cleanup:
```go
func BenchmarkProcess(b *testing.B) {
    setup := prepareTestData()

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        b.StopTimer()
        data := cloneTestData(setup) // expensive per-iteration setup
        b.StartTimer()

        Process(data)
    }

    // Optional cleanup after benchmark
    cleanupTestData(setup)
}
```
Template for new benchmarks:
```go
func BenchmarkXxx(b *testing.B) {
    // Global setup
    setup := InitSetup()

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        b.StopTimer()
        item := PrepareData(setup)
        b.StartTimer()

        TargetFunction(item)
    }

    Cleanup(setup)
}
```

⸻

## 🏃 Useful Commands

Run unit tests:
```bash 
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run benchmarks:
```bash
go test -bench . ./...
```

Run the project:
```bash
go run main.go
```
Build and run the project:
```bash
go build -o csv2parquet main.go
./csv2parquet
```

⸻

## ✅ Summary

Use this document to guide your completions. Ensure that your code:
Ensure that your code:
-	Is idiomatic and clean
-	Follows project structure
-	Matches naming and formatting standards
-	Has full unit test coverage with proper cases
-	Provides accurate and efficient benchmarks
-	Supports project running and testing commands
