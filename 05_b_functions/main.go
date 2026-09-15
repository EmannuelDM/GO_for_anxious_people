package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 1. ERROR HANDLING & SENTINEL ERRORS
// ============================================================================

// Sentinel errors: package-level variables representing specific, known error conditions.
// By convention, sentinel errors start with 'Err'.
var (
	ErrNotFound       = errors.New("record not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidConfig  = errors.New("invalid configuration")
)

// Custom error struct implementing the built-in 'error' interface:
// type error interface { Error() string }
type ValidationError struct {
	Field   string
	Value   any
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s' (value: %v): %s", e.Field, e.Value, e.Message)
}

// simulateDatabaseQuery simulates a lookup that can fail with sentinel errors.
func simulateDatabaseQuery(id int) (string, error) {
	if id < 0 {
		return "", ErrPermissionDenied
	}
	if id == 0 {
		return "", ErrNotFound
	}
	return fmt.Sprintf("User#%d", id), nil
}

// fetchUserData demonstrates ERROR WRAPPING using '%w' in fmt.Errorf.
// Wrapping preserves the error chain so callers can inspect underlying causes with errors.Is/errors.As.
func fetchUserData(id int) (string, error) {
	data, err := simulateDatabaseQuery(id)
	if err != nil {
		// '%w' wraps the original error; '%v' would only format its string without preserving identity
		return "", fmt.Errorf("fetchUserData(id=%d): %w", id, err)
	}
	return data, nil
}

// validateUserForm demonstrates returning custom error types and joining multiple errors.
func validateUserForm(username string, age int) error {
	var errs []error

	if strings.TrimSpace(username) == "" {
		errs = append(errs, &ValidationError{
			Field:   "username",
			Value:   username,
			Message: "cannot be empty",
		})
	}

	if age < 18 || age > 120 {
		errs = append(errs, &ValidationError{
			Field:   "age",
			Value:   age,
			Message: "must be between 18 and 120",
		})
	}

	// errors.Join (Go 1.20+) combines multiple errors into a single error.
	// If the slice is empty, errors.Join returns nil!
	return errors.Join(errs...)
}

// ============================================================================
// 2. WORKING WITH EOF (io.EOF)
// ============================================================================

// demonstrateLowLevelEOF demonstrates the crucial io.Reader contract:
// When Read() reaches the end of stream, it may return n > 0 bytes AND err == io.EOF
// at the same time. ALWAYS process n bytes first before checking err == io.EOF!
func demonstrateLowLevelEOF(streamData string) error {
	reader := strings.NewReader(streamData)
	buf := make([]byte, 8) // Small 8-byte buffer to force multiple read cycles

	totalBytesRead := 0
	readCycle := 0

	for {
		readCycle++
		n, err := reader.Read(buf)

		// STEP 1: Process valid data read in this cycle (even if err != nil!)
		if n > 0 {
			totalBytesRead += n
			fmt.Printf("    Cycle %d: read %d bytes -> %q\n", readCycle, n, string(buf[:n]))
		}

		// STEP 2: Check for error or EOF
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("    Cycle %d: reached io.EOF gracefully! Total bytes: %d\n", readCycle, totalBytesRead)
				return nil // Graceful completion, EOF is not a failure!
			}
			// Any error other than EOF is a real, unexpected failure
			return fmt.Errorf("unexpected read error: %w", err)
		}
	}
}

// demonstrateBufioReaderEOF demonstrates line-by-line reading where the last line
// might not have a trailing newline before EOF.
func demonstrateBufioReaderEOF(content string) {
	reader := bufio.NewReader(strings.NewReader(content))
	lineNumber := 0

	for {
		lineNumber++
		line, err := reader.ReadString('\n')

		// Even when err == io.EOF, ReadString may return data if the last line lacked '\n'
		if len(line) > 0 {
			cleanLine := strings.TrimRight(line, "\r\n")
			fmt.Printf("    Line %d: %s\n", lineNumber, cleanLine)
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("    [bufio.Reader] Reached io.EOF, stopping cleanly.")
				break
			}
			fmt.Printf("    [bufio.Reader] Error: %v\n", err)
			break
		}
	}
}

// demonstrateScannerEOF demonstrates how bufio.Scanner abstracts away io.EOF internally.
func demonstrateScannerEOF(content string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	count := 0

	// scanner.Scan() automatically consumes io.EOF and returns false when finished!
	for scanner.Scan() {
		count++
		fmt.Printf("    Token %d: %s\n", count, scanner.Text())
	}

	// scanner.Err() returns any non-EOF error. If EOF caused Scan to stop, Err() returns nil.
	if err := scanner.Err(); err != nil {
		fmt.Printf("    Scanner error: %v\n", err)
	} else {
		fmt.Println("    [bufio.Scanner] Scan finished with no errors (io.EOF handled internally).")
	}
}

// ============================================================================
// 3. ANONYMOUS FUNCTIONS & CLOSURES
// ============================================================================

// makeAccumulator returns a closure that maintains state between calls.
// The inner function captures 'sum' from its enclosing lexical environment.
func makeAccumulator(initial int) func(int) int {
	sum := initial
	return func(val int) int {
		sum += val
		return sum
	}
}

// makeRateLimiter creates a closure that limits actions per token bucket.
func makeRateLimiter(maxRequests int) func() bool {
	tokens := maxRequests
	return func() bool {
		if tokens > 0 {
			tokens--
			return true // Allowed
		}
		return false // Rate limited
	}
}

// demonstrateAnonymousFunctions explores various anonymous function patterns:
// 1. Immediately Invoked Function Expressions (IIFE)
// 2. Assigning to variables
// 3. Inline callbacks for sorting/filtering
// 4. Deferred closures modifying named returns or measuring duration
func demonstrateAnonymousFunctions() {
	// Pattern 1: IIFE (Immediately Invoked Function Expression)
	// Declared and executed in a single line, useful for isolating temporary setup scope.
	configVersion := func(major, minor int) string {
		return fmt.Sprintf("v%d.%d-release", major, minor)
	}(2, 5)
	fmt.Printf("  IIFE result: %s\n", configVersion)

	// Pattern 2: Assigning anonymous function to variable
	sanitize := func(input string) string {
		return strings.ToLower(strings.TrimSpace(input))
	}
	fmt.Printf("  Sanitize variable call: %q\n", sanitize("   GoLang Rocks!   "))

	// Pattern 3: Inline callback in standard library (e.g. sort.Slice)
	words := []string{"elephant", "cat", "hippopotamus", "dog"}
	fmt.Printf("  Original words: %v\n", words)

	// Passing anonymous predicate function directly
	sort.Slice(words, func(i, j int) bool {
		return len(words[i]) < len(words[j]) // Sort ascending by length
	})
	fmt.Printf("  Sorted by length: %v\n", words)

	// Pattern 4: Stateful Closures
	acc := makeAccumulator(100)
	fmt.Printf("  Accumulator initial 100 + 50 = %d\n", acc(50))
	fmt.Printf("  Accumulator + 25 = %d\n", acc(25))
	fmt.Printf("  Accumulator + 10 = %d\n", acc(10))

	limiter := makeRateLimiter(2)
	fmt.Printf("  Limiter request 1: %t (allowed)\n", limiter())
	fmt.Printf("  Limiter request 2: %t (allowed)\n", limiter())
	fmt.Printf("  Limiter request 3: %t (blocked!)\n", limiter())

	// Pattern 5: Anonymous function in defer for timing & telemetry
	func() {
		defer func(start time.Time) {
			fmt.Printf("  [telemetry closure] Operation completed in %v\n", time.Since(start))
		}(time.Now())

		// Simulated short work
		time.Sleep(15 * time.Millisecond)
	}()

	// Pattern 6: Anonymous functions in goroutines (safe parameter passing)
	var wg sync.WaitGroup
	items := []string{"alpha", "beta", "gamma"}

	fmt.Println("  Launching concurrent workers via anonymous goroutines:")
	for _, item := range items {
		wg.Add(1)
		// Passing 'item' as argument to the anonymous function avoids closure data races
		go func(val string) {
			defer wg.Done()
			fmt.Printf("    Worker processing: %s\n", val)
		}(item)
	}
	wg.Wait()
}

// ============================================================================
// 4. VARIADIC FUNCTIONS (...T)
// ============================================================================

// Variadic function accepting variable number of integers.
// Inside the function, 'numbers' has the type []int.
func calculateMetrics(category string, numbers ...int) (min int, max int, sum int, count int) {
	count = len(numbers)
	if count == 0 {
		return 0, 0, 0, 0
	}

	min = numbers[0]
	max = numbers[0]
	sum = 0

	for _, n := range numbers {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
		sum += n
	}
	return min, max, sum, count
}

// Variadic function with 'any' (...any) for type-flexible formatting or logging.
// Similar to how fmt.Println is implemented.
func formatSummary(prefix string, values ...any) string {
	var parts []string
	for _, v := range values {
		switch val := v.(type) {
		case string:
			parts = append(parts, fmt.Sprintf("%q (string)", val))
		case int:
			parts = append(parts, fmt.Sprintf("%d (int)", val))
		case bool:
			parts = append(parts, fmt.Sprintf("%t (bool)", val))
		case error:
			parts = append(parts, fmt.Sprintf("[%v] (error)", val))
		default:
			parts = append(parts, fmt.Sprintf("%v (%T)", val, val))
		}
	}
	return fmt.Sprintf("[%s] %s", prefix, strings.Join(parts, ", "))
}

// ============================================================================
// 5. THE REAL-WORLD SYNERGY: FUNCTIONAL OPTIONS PATTERN
// ============================================================================
// This industry-standard pattern combines:
// 1. Variadic functions (opts ...ServerOption)
// 2. Anonymous functions / closures (func(*ServerConfig) error)
// 3. Error handling & validation

type ServerConfig struct {
	Host         string
	Port         int
	Timeout      time.Duration
	EnableTLS    bool
	MaxClients   int
}

// ServerOption is a function type that mutates ServerConfig and can return errors.
type ServerOption func(*ServerConfig) error

// Option constructors returning closures:

func WithPort(port int) ServerOption {
	return func(s *ServerConfig) error {
		if port < 1 || port > 65535 {
			return fmt.Errorf("invalid port %d: %w", port, ErrInvalidConfig)
		}
		s.Port = port
		return nil
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *ServerConfig) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must be positive (got %v): %w", timeout, ErrInvalidConfig)
		}
		s.Timeout = timeout
		return nil
	}
}

func WithTLS(enabled bool) ServerOption {
	return func(s *ServerConfig) error {
		s.EnableTLS = enabled
		return nil
	}
}

func WithMaxClients(maxClients int) ServerOption {
	return func(s *ServerConfig) error {
		if maxClients <= 0 {
			return fmt.Errorf("maxClients must be > 0: %w", ErrInvalidConfig)
		}
		s.MaxClients = maxClients
		return nil
	}
}

// NewServerConfig uses variadic options to configure a struct safely with default fallbacks.
func NewServerConfig(host string, opts ...ServerOption) (*ServerConfig, error) {
	if strings.TrimSpace(host) == "" {
		return nil, errors.New("host cannot be empty")
	}

	// Default configuration values
	cfg := &ServerConfig{
		Host:       host,
		Port:       8080,
		Timeout:    30 * time.Second,
		EnableTLS:  false,
		MaxClients: 100,
	}

	// Apply each variadic option (which are closures executing on cfg)
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("failed applying server option: %w", err)
		}
	}

	return cfg, nil
}

// ============================================================================
// MAIN FUNCTION & DEMO RUNNERS
// ============================================================================

func main() {
	fmt.Println("================================================================")
	fmt.Println("   05_b: ERRORS, EOF, ANONYMOUS FUNCTIONS & VARIADIC FUNCTIONS  ")
	fmt.Println("================================================================")

	demoErrorHandling()
	demoEOFHandling()
	demoAnonymousFunctions()
	demoVariadicFunctions()
	demoFunctionalOptionsPattern()
}

// ----------------------------------------------------------------------------
// Demo 1: Errors, Sentinel Errors, errors.Is, errors.As, errors.Join
// ----------------------------------------------------------------------------
func demoErrorHandling() {
	fmt.Println("\n--- 1. Errors: Sentinel, Wrapping (%w), errors.Is & errors.As ---")

	// Testing sentinel error inspection with errors.Is
	_, errNotFound := fetchUserData(0)
	if errNotFound != nil {
		fmt.Printf("Caught error: %v\n", errNotFound)

		// errors.Is recursively unwraps the error chain created with %w!
		if errors.Is(errNotFound, ErrNotFound) {
			fmt.Println("  -> errors.Is confirmed: The root cause was ErrNotFound!")
		}
	}

	_, errDenied := fetchUserData(-5)
	if errDenied != nil {
		fmt.Printf("Caught error: %v\n", errDenied)
		if errors.Is(errDenied, ErrPermissionDenied) {
			fmt.Println("  -> errors.Is confirmed: Access was ErrPermissionDenied!")
		}
	}

	// Testing custom error structs with errors.As and errors.Join
	fmt.Println("\nValidating user form (collecting multiple errors):")
	validationErr := validateUserForm("", 15)
	if validationErr != nil {
		fmt.Println("  Combined error output:")
		fmt.Printf("    %v\n", validationErr)

		// Using errors.As to extract structured fields from the first matching ValidationError
		var targetErr *ValidationError
		if errors.As(validationErr, &targetErr) {
			fmt.Printf("  -> errors.As extracted typed struct: Field=%q, Reason=%q\n",
				targetErr.Field, targetErr.Message)
		}
	}
}

// ----------------------------------------------------------------------------
// Demo 2: Handling EOF (io.EOF)
// ----------------------------------------------------------------------------
func demoEOFHandling() {
	fmt.Println("\n--- 2. Working with EOF (io.EOF) ---")

	sampleText := "Hello, Go Stream EOF mechanics!"

	// 2.1: Low-level chunked reading
	fmt.Println("2.1 Low-level Reader (demonstrating n > 0 before checking io.EOF):")
	if err := demonstrateLowLevelEOF(sampleText); err != nil {
		fmt.Printf("Reader failed: %v\n", err)
	}

	// 2.2: bufio.Reader line-by-line (handling last line without newline)
	fmt.Println("\n2.2 bufio.Reader (ReadString until EOF):")
	multilineContent := "Row 1: Alpha\nRow 2: Beta\nRow 3: Gamma (no newline at end)"
	demonstrateBufioReaderEOF(multilineContent)

	// 2.3: bufio.Scanner (abstracting io.EOF internally)
	fmt.Println("\n2.3 bufio.Scanner (EOF consumed internally):")
	csvContent := "item_apple\nitem_banana\nitem_cherry"
	demonstrateScannerEOF(csvContent)
}

// ----------------------------------------------------------------------------
// Demo 3: Anonymous Functions and Closures
// ----------------------------------------------------------------------------
func demoAnonymousFunctions() {
	fmt.Println("\n--- 3. Anonymous Functions & Closures ---")
	demonstrateAnonymousFunctions()
}

// ----------------------------------------------------------------------------
// Demo 4: Variadic Functions
// ----------------------------------------------------------------------------
func demoVariadicFunctions() {
	fmt.Println("\n--- 4. Variadic Functions (...T and ...any) ---")

	// Calling with 0 arguments
	min, max, sum, count := calculateMetrics("Empty")
	fmt.Printf("  Empty call -> count: %d, sum: %d\n", count, sum)

	// Calling with individual discrete arguments
	min, max, sum, count = calculateMetrics("Temperatures", 18, 24, 15, 30, 22)
	fmt.Printf("  Individual arguments -> count: %d, min: %d, max: %d, sum: %d\n", count, min, max, sum)

	// Unpacking a slice with spread operator '...'
	measurements := []int{105, 98, 112, 120}
	min, max, sum, count = calculateMetrics("Slice Unpacked", measurements...)
	fmt.Printf("  Unpacked slice (%v) -> count: %d, min: %d, max: %d, sum: %d\n",
		measurements, count, min, max, sum)

	// Built-in append is also variadic!
	baseSlice := []int{1, 2, 3}
	extraSlice := []int{4, 5, 6}
	combined := append(baseSlice, extraSlice...) // Unpacking slice into append
	fmt.Printf("  Built-in append with spread: %v\n", combined)

	// Variadic with any (...any) and type inspection
	logOutput := formatSummary("AuditLog", "User logged in", 404, false, ErrNotFound, 3.14159)
	fmt.Printf("  formatSummary(...any): %s\n", logOutput)
}

// ----------------------------------------------------------------------------
// Demo 5: Real-World Synergy: Functional Options Pattern
// ----------------------------------------------------------------------------
func demoFunctionalOptionsPattern() {
	fmt.Println("\n--- 5. Real-World Synergy: Functional Options Pattern ---")

	// Case 1: Default configuration (zero variadic options)
	defaultSrv, err := NewServerConfig("api.example.com")
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("  Default Server: Host=%s, Port=%d, Timeout=%v, TLS=%t, MaxClients=%d\n",
			defaultSrv.Host, defaultSrv.Port, defaultSrv.Timeout, defaultSrv.EnableTLS, defaultSrv.MaxClients)
	}

	// Case 2: Custom configuration using variadic options
	customSrv, err := NewServerConfig("secure.payment.gateway",
		WithPort(443),
		WithTLS(true),
		WithTimeout(10*time.Second),
		WithMaxClients(500),
	)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("  Custom Server:  Host=%s, Port=%d, Timeout=%v, TLS=%t, MaxClients=%d\n",
			customSrv.Host, customSrv.Port, customSrv.Timeout, customSrv.EnableTLS, customSrv.MaxClients)
	}

	// Case 3: Error handling within functional options
	_, errInvalid := NewServerConfig("bad.config.com",
		WithPort(99999), // Exceeds valid port range!
	)
	if errInvalid != nil {
		fmt.Printf("  Invalid option rejected: %v\n", errInvalid)
		if errors.Is(errInvalid, ErrInvalidConfig) {
			fmt.Println("  -> Confirmed error is ErrInvalidConfig via errors.Is!")
		}
	}
}
