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

type Config struct {
	CSVFile      string
	TemplateFile string
	OutputDir    string
}

type RowData map[string]string

func main() {
	config := parseFlags()
	validateInput(config)
	records, headers := readCSV(config.CSVFile)
	templateContent := readTemplateFile(config.TemplateFile)
	createOutputDir(config.OutputDir)
	processRows(records, headers, templateContent, config)
}

func parseFlags() Config {
	csvFile := flag.String("csv", "", "Location of the CSV file")
	templateFile := flag.String("template", "", "Location of the template file")
	outputDir := flag.String("output", "", "Location of the output directory")
	flag.Parse()
	return Config{*csvFile, *templateFile, *outputDir}
}

func validateInput(config Config) {
	if config.CSVFile == "" || config.TemplateFile == "" || config.OutputDir == "" {
		log.Fatal("All arguments (csv, template, output) are required")
	}
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

func readTemplateFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading template file: %v", err)
	}
	return string(content)
}

func createOutputDir(outputDir string) {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}
}

func processRows(records [][]string, headers []string, templateContent string, config Config) {
	var wg sync.WaitGroup
	for i, record := range records {
		wg.Add(1)
		go func(rowNum int, row []string) {
			defer wg.Done()
			processRow(rowNum, row, headers, templateContent, config)
		}(i+1, record)
	}
	wg.Wait()
}

func processRow(rowNum int, row []string, headers []string, templateContent string, config Config) {
	data := makeRowData(row, headers)
	processedContent := applyTemplate(templateContent, data)
	writeOutputFile(rowNum, processedContent, config.TemplateFile, config.OutputDir)
}

func makeRowData(row []string, headers []string) RowData {
	data := make(RowData)
	for i, value := range row {
		data[headers[i]] = value
	}
	return data
}

func applyTemplate(templateContent string, data RowData) string {
	tmpl, err := template.New("row").Parse(templateContent)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	var processed strings.Builder
	err = tmpl.Execute(&processed, data)
	if err != nil {
		log.Fatalf("Error applying template: %v", err)
	}

	return processed.String()
}

func writeOutputFile(rowNum int, content, templateFile, outputDir string) {
	ext := filepath.Ext(templateFile)
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%d%s", rowNum, ext))

	err := ioutil.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		log.Printf("Error writing output file for row %d: %v", rowNum, err)
		return
	}

	fmt.Printf("Processed row %d, output: %s\n", rowNum, outputFile)
}