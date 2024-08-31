# CSV Template Tool

A command-line tool that processes CSV data and applies it to a template file, generating output files for each CSV row.

## Getting the Tool

```bash
git clone https://github.com/patelrohanv/csv-template-tool.git
cd csv-template-tool
go build
```

## Usage

```bash
./csv-template-tool -csv=<csv_file> -template=<template_file> -output=<output_directory>
```

Flags:
- `-csv`: Path to the input CSV file
- `-template`: Path to the template file
- `-output`: Path to the output directory

## Running with Sample Files

```bash
./csv-template-tool -csv=sample.csv -template=sample.json -output=output
```

This command reads from `sample.csv`, applies the data to `sample.json`, and generates output files in the `./output` directory.

## Error Checking

Check the console output for any error messages or warnings during execution.