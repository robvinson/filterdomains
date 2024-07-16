package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: filterdomains -fd <path to filter file>")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	// Parse command-line arguments
	filterFilePath := flag.String("fd", "", "Path to the file containing domains to filter out")
	flag.Parse()

	if *filterFilePath == "" {
		usage()
	}

	// Read the filter file
	filterFile, err := os.Open(*filterFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening filter file: %v\n", err)
		os.Exit(1)
	}
	defer filterFile.Close()

	// Read filter domains into a slice
	var filters []string
	scanner := bufio.NewScanner(filterFile)
	for scanner.Scan() {
		filters = append(filters, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading filter file: %v\n", err)
		os.Exit(1)
	}

	// Create a scanner to read from stdin
	inputScanner := bufio.NewScanner(os.Stdin)

	// Iterate over each line (domain name) from stdin
	for inputScanner.Scan() {
		domain := strings.TrimSpace(inputScanner.Text())
		shouldFilter := false

		// Check if the domain matches any of the filter patterns
		for _, pattern := range filters {
			matched, err := filepath.Match(pattern, domain)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error matching pattern: %v\n", err)
				os.Exit(1)
			}
			if matched {
				shouldFilter = true
				break
			}
		}

		// Print the domain if it does not match any filter patterns
		if !shouldFilter {
			fmt.Println(domain)
		}
	}

	// Check for scanner errors
	if err := inputScanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		os.Exit(1)
	}
}
