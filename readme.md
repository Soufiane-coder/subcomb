# SubPerm - Subdomain Permutation Generator

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

A powerful and flexible command-line tool written in Go that generates all possible subdomain permutations from input subdomains. Perfect for security testing, domain enumeration, and reconnaissance tasks.

## Features

- 🚀 **Fast and Efficient** - Written in Go for optimal performance
- 📁 **File I/O Support** - Read from files and write to files
- 🔄 **Pipeline Support** - Works seamlessly with Unix pipelines
- 📊 **Multiple Output Formats** - Plain text, JSON, and CSV
- 🎯 **Smart Domain Detection** - Validates and processes only valid subdomains
- 🔧 **Flexible Options** - Comprehensive command-line interface
- 🧹 **Duplicate Removal** - Automatic deduplication of results

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/Soufiane-coder/subperm.git
cd subperm

# Build the binary
go build -o subperm main.go

# Make it executable and install (optional)
chmod +x subperm
sudo mv subperm /usr/local/bin/
```

### Direct Download

```bash
# Download and install directly
go install github.com/Soufiane-coder/subperm@latest
```

## Usage

### Basic Usage

```bash
# Generate permutations for a single subdomain
subperm sub.api.example.com
```

**Output:**
```
example.com
sub.example.com
api.example.com
api.sub.example.com
sub.api.example.com
```

### Command Line Options

```
USAGE:
    subperm [OPTIONS] [SUBDOMAIN]

OPTIONS:
    -i, --input FILE     Read subdomains from input file
    -o, --output FILE    Write results to output file (default: stdout)
    -u, --unique         Remove duplicate results (default: true)
    -v, --verbose        Enable verbose output
    -h, --help           Show this help message
    --version            Show version information
    -f, --format FORMAT  Output format: plain, json, csv (default: plain)
```

### Examples

#### File Operations

```bash
# Read from file and write to file
subperm -i subdomains.txt -o results.txt

# Verbose mode with file processing
subperm -v -i input.txt -o output.txt
```

#### Pipeline Operations

```bash
# Pipeline from echo
echo "sub.api.example.com" | subperm

# Pipeline from file
cat subdomains.txt | subperm

# Chain with other Unix tools
subperm -i domains.txt | grep "api" | sort | uniq

# Advanced pipeline usage
cat targets.txt | subperm | httpx -silent | nuclei -t vulnerabilities/
```

#### Output Formats

```bash
# JSON output
subperm -f json sub.api.example.com
# Output: ["example.com","sub.example.com","api.example.com","api.sub.example.com","sub.api.example.com"]

# CSV output
subperm -f csv -i input.txt -o results.csv

# Plain text (default)
subperm sub.api.example.com > results.txt
```

## Input File Format

Create a text file with one subdomain per line:

```
# subdomains.txt
sub.api.example.com
test.staging.mysite.org
dev.app.company.com
admin.panel.site.net

# Comments are supported and ignored
# Empty lines are also ignored

www.blog.example.com
```

## Output Examples

### Plain Text Format
```
example.com
sub.example.com
api.example.com
staging.example.com
api.sub.example.com
staging.sub.example.com
sub.api.example.com
sub.staging.example.com
api.staging.example.com
staging.api.example.com
api.staging.sub.example.com
staging.api.sub.example.com
sub.api.staging.example.com
sub.staging.api.example.com
api.sub.staging.example.com
staging.sub.api.example.com
```

### JSON Format
```json
[
  "example.com",
  "sub.example.com",
  "api.example.com",
  "staging.example.com",
  "api.sub.example.com",
  "staging.sub.example.com"
]
```

### CSV Format
```csv
subdomain
example.com
sub.example.com
api.example.com
staging.example.com
api.sub.example.com
staging.sub.example.com
```

## Use Cases

### Security Testing
```bash
# Generate subdomains for reconnaissance
subperm -i target-domains.txt | httpx -silent -title -tech-detect

# Find live subdomains
cat subdomains.txt | subperm | httprobe | tee live-domains.txt
```

### Domain Enumeration
```bash
# Comprehensive subdomain discovery
subperm -i seeds.txt | dnsx -silent -a -resp | anew discovered-domains.txt
```

### Integration with Other Tools
```bash
# With Subfinder and Httpx
subfinder -d example.com -silent | subperm | httpx -silent -mc 200

# With Amass
amass enum -d example.com | subperm -f json > permutations.json

# With custom scripts
subperm -i domains.txt | while read domain; do
    dig +short "$domain" | grep -E '^[0-9]'
done
```

## Algorithm

The tool works by:

1. **Parsing** the input subdomain into components
2. **Extracting** the base domain and TLD
3. **Identifying** subdomain parts
4. **Generating** all possible permutations and combinations
5. **Removing** duplicates (if enabled)
6. **Outputting** in the specified format

For `sub.api.example.com`, it generates:
- Base domain: `example.com`
- Single subdomains: `sub.example.com`, `api.example.com`
- Combined subdomains: `api.sub.example.com`, `sub.api.example.com`

## Performance

- **Memory Efficient**: Processes large files without loading everything into memory
- **Fast Processing**: Optimized algorithms for permutation generation
- **Scalable**: Handles thousands of subdomains efficiently

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/Soufiane-coder/subperm.git
cd subperm

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build
go build -o subperm main.go
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

### v1.0.0
- Initial release
- Basic subdomain permutation generation
- File I/O support
- Pipeline support
- Multiple output formats
- Comprehensive CLI interface

## Support

If you find this tool useful, please consider:
- ⭐ Starring the repository
- 🐛 Reporting bugs
- 💡 Suggesting new features
- 🤝 Contributing code

## Acknowledgments

- Inspired by various subdomain enumeration tools in the security community
- Built with the Go standard library for maximum compatibility
- Designed to integrate seamlessly with existing security workflows

---

**Made with ❤️ for the cybersecurity community**