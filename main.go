package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

type RowData map[string]string

func main() {
	// Define command-line flags
	csvFile := flag.String("csv", "", "Location of the CSV file")
	templateFile := flag.String("template", "", "Location of the template file")
	outputDir := flag.String("output", "", "Location of the output directory")
	flag.Parse()

	// Validate input
	if *csvFile == "" || *templateFile == "" || *outputDir == "" {
		log.Fatal("All arguments (csv, template, output) are required")
	}

	// Read CSV file
	records, headers := readCSV(*csvFile)

	// Read template file
	templateContent, err := ioutil.ReadFile(*templateFile)
	if err != nil {
		log.Fatalf("Error reading template file: %v", err)
	}

	// Create output directory if it doesn't exist
	err = os.MkdirAll(*outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	// Process rows concurrently
	var wg sync.WaitGroup
	for i, record := range records {
		wg.Add(1)
		go func(rowNum int, row []string) {
			defer wg.Done()
			processRow(rowNum, row, headers, string(templateContent), *templateFile, *outputDir)
		}(i+1, record)
	}
	wg.Wait()

	fmt.Println("Processing complete.")
}

func readCSV(filename string) ([][]string, []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		log.Fatal("CSV file must have at least a header row and one data row")
	}

	return records[1:], records[0]
}

func processRow(rowNum int, row []string, headers []string, templateContent, templateFile, outputDir string) {
	// Create a map of the row data
	data := make(RowData)
	for i, value := range row {
		data[headers[i]] = value
	}

	// Parse the template
	tmpl, err := template.New("row").Parse(templateContent)
	if err != nil {
		log.Printf("Error parsing template for row %d: %v", rowNum, err)
		return
	}

	// Apply the template
	var processed strings.Builder
	err = tmpl.Execute(&processed, data)
	if err != nil {
		log.Printf("Error applying template for row %d: %v", rowNum, err)
		return
	}

	// Determine output file name
	ext := filepath.Ext(templateFile)
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%d%s", rowNum, ext))

	// Write the processed content to the output file
	err = ioutil.WriteFile(outputFile, []byte(processed.String()), 0644)
	if err != nil {
		log.Printf("Error writing output file for row %d: %v", rowNum, err)
		return
	}

	fmt.Printf("Processed row %d, output: %s\n", rowNum, outputFile)
}
