package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	version = "1.0.0"
	usage   = `subcomb - Subdomain Permutation Generator

USAGE:
    subcomb [OPTIONS] [SUBDOMAIN]

OPTIONS:
    -i, --input FILE     Read subdomains from input file
    -o, --output FILE    Write results to output file (default: stdout)
    -u, --unique         Remove duplicate results (default: true)
    -v, --verbose        Enable verbose output
    -h, --help           Show this help message
    --version            Show version information
    -f, --format FORMAT  Output format: plain, json, csv (default: plain)

EXAMPLES:
    # Single subdomain
    subcomb sub.api.example.com

    # Read from file
    subcomb -i subdomains.txt -o results.txt

    # Pipeline usage
    echo "sub.api.example.com" | subcomb
    cat subdomains.txt | subcomb -o results.txt

    # Multiple formats
    subcomb -f json sub.api.example.com
    subcomb -f csv -i input.txt -o output.csv

FORMATS:
    plain  - One subdomain per line (default)
    json   - JSON array format
    csv    - Comma-separated values with header
`
)

type Config struct {
	InputFile   string
	OutputFile  string
	Unique      bool
	Verbose     bool
	Format      string
	ShowHelp    bool
	ShowVersion bool
}

type SubdomainGenerator struct {
	config *Config
}

// NewSubdomainGenerator creates a new generator with config
func NewSubdomainGenerator(config *Config) *SubdomainGenerator {
	return &SubdomainGenerator{config: config}
}

// ParseFlags parses command line flags
func ParseFlags() (*Config, []string) {
	config := &Config{}

	flag.StringVar(&config.InputFile, "i", "", "Input file")
	flag.StringVar(&config.InputFile, "input", "", "Input file")
	flag.StringVar(&config.OutputFile, "o", "", "Output file")
	flag.StringVar(&config.OutputFile, "output", "", "Output file")
	flag.BoolVar(&config.Unique, "u", true, "Remove duplicates")
	flag.BoolVar(&config.Unique, "unique", true, "Remove duplicates")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose output")
	flag.BoolVar(&config.Verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&config.ShowHelp, "h", false, "Show help")
	flag.BoolVar(&config.ShowHelp, "help", false, "Show help")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version")
	flag.StringVar(&config.Format, "f", "plain", "Output format")
	flag.StringVar(&config.Format, "format", "plain", "Output format")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()

	return config, flag.Args()
}

// IsValidSubdomain checks if a string looks like a subdomain
func (sg *SubdomainGenerator) IsValidSubdomain(input string) bool {
	// Simple regex to match domain-like patterns
	pattern := `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, input)
	return err == nil && matched
}

// GeneratePermutations creates all possible subdomain combinations
func (sg *SubdomainGenerator) GeneratePermutations(subdomain string) []string {
	subdomain = strings.TrimSpace(strings.ToLower(subdomain))
	subdomain = strings.TrimSuffix(subdomain, ".")

	if !sg.IsValidSubdomain(subdomain) {
		if sg.config.Verbose {
			fmt.Fprintf(os.Stderr, "Warning: '%s' doesn't appear to be a valid subdomain\n", subdomain)
		}
		return []string{}
	}

	parts := strings.Split(subdomain, ".")
	if len(parts) < 2 {
		return []string{subdomain}
	}

	// Extract domain and TLD (last two parts)
	domain := parts[len(parts)-2]
	tld := parts[len(parts)-1]
	subdomainParts := parts[:len(parts)-2]

	var combinations []string
	baseDomain := domain + "." + tld

	// Add base domain
	combinations = append(combinations, baseDomain)

	// Generate all combinations
	if len(subdomainParts) > 0 {
		combinations = append(combinations, sg.generateCombinations(subdomainParts, baseDomain)...)
	}

	if sg.config.Unique {
		combinations = sg.removeDuplicates(combinations)
	}

	return combinations
}

// generateCombinations creates all permutations of subdomain parts
func (sg *SubdomainGenerator) generateCombinations(parts []string, baseDomain string) []string {
	var result []string

	// Generate all possible lengths from 1 to len(parts)
	for length := 1; length <= len(parts); length++ {
		sg.generatePermutationsOfLength(parts, length, []string{}, &result, baseDomain)
	}

	return result
}

// generatePermutationsOfLength generates all permutations of specific length
func (sg *SubdomainGenerator) generatePermutationsOfLength(parts []string, length int, current []string, result *[]string, baseDomain string) {
	if len(current) == length {
		subdomain := strings.Join(current, ".") + "." + baseDomain
		*result = append(*result, subdomain)
		return
	}

	for i := 0; i < len(parts); i++ {
		if !sg.contains(current, parts[i]) {
			newCurrent := make([]string, len(current)+1)
			copy(newCurrent, current)
			newCurrent[len(current)] = parts[i]
			sg.generatePermutationsOfLength(parts, length, newCurrent, result, baseDomain)
		}
	}
}

// contains checks if slice contains item
func (sg *SubdomainGenerator) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// removeDuplicates removes duplicate entries
func (sg *SubdomainGenerator) removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

// ProcessInput reads from input source and processes subdomains
func (sg *SubdomainGenerator) ProcessInput(reader io.Reader) ([]string, error) {
	var allResults []string
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		if sg.config.Verbose {
			fmt.Fprintf(os.Stderr, "Processing: %s\n", line)
		}

		results := sg.GeneratePermutations(line)
		allResults = append(allResults, results...)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	if sg.config.Unique {
		allResults = sg.removeDuplicates(allResults)
	}

	return allResults, nil
}

// WriteOutput writes results to output destination
func (sg *SubdomainGenerator) WriteOutput(writer io.Writer, results []string) error {
	switch sg.config.Format {
	case "json":
		return sg.writeJSON(writer, results)
	case "csv":
		return sg.writeCSV(writer, results)
	default:
		return sg.writePlain(writer, results)
	}
}

// writePlain writes plain text output
func (sg *SubdomainGenerator) writePlain(writer io.Writer, results []string) error {
	for _, result := range results {
		if _, err := fmt.Fprintln(writer, result); err != nil {
			return err
		}
	}
	return nil
}

// writeJSON writes JSON output
func (sg *SubdomainGenerator) writeJSON(writer io.Writer, results []string) error {
	fmt.Fprint(writer, "[")
	for i, result := range results {
		if i > 0 {
			fmt.Fprint(writer, ",")
		}
		fmt.Fprintf(writer, "\"%s\"", result)
	}
	fmt.Fprintln(writer, "]")
	return nil
}

// writeCSV writes CSV output
func (sg *SubdomainGenerator) writeCSV(writer io.Writer, results []string) error {
	fmt.Fprintln(writer, "subdomain")
	for _, result := range results {
		fmt.Fprintln(writer, result)
	}
	return nil
}

// Run executes the main logic
func (sg *SubdomainGenerator) Run(args []string) error {
	var reader io.Reader
	var writer io.Writer = os.Stdout

	// Determine input source
	if sg.config.InputFile != "" {
		file, err := os.Open(sg.config.InputFile)
		if err != nil {
			return fmt.Errorf("error opening input file: %v", err)
		}
		defer file.Close()
		reader = file
	} else if len(args) > 0 {
		// Use command line argument
		reader = strings.NewReader(args[0])
	} else {
		// Check if stdin has data (pipeline usage)
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			reader = os.Stdin
		} else {
			return fmt.Errorf("no input provided. Use -h for help")
		}
	}

	// Determine output destination
	if sg.config.OutputFile != "" {
		file, err := os.Create(sg.config.OutputFile)
		if err != nil {
			return fmt.Errorf("error creating output file: %v", err)
		}
		defer file.Close()
		writer = file
	}

	// Process input
	results, err := sg.ProcessInput(reader)
	if err != nil {
		return err
	}

	if sg.config.Verbose {
		fmt.Fprintf(os.Stderr, "Generated %d results\n", len(results))
	}

	// Write output
	return sg.WriteOutput(writer, results)
}

func main() {
	config, args := ParseFlags()

	if config.ShowVersion {
		fmt.Printf("subcomb version %s\n", version)
		return
	}

	if config.ShowHelp {
		fmt.Print(usage)
		return
	}

	generator := NewSubdomainGenerator(config)

	if err := generator.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}