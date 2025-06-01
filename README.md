# openfoodfacts-to-eatnlift

## About The Project

This project is provided in compliance with Section 4.6 b. of the [ODbL](https://opendatacommons.org/licenses/odbl/1-0/).

Running this project converts a single Open Food Facts JSONL gzipped Data Export into a set of JSONL files. Apart from file type conversion and chunking into multiple smaller files, this project also mutates the data structure, transforms some of the values and drops certain rows and columns as per the needs of the Eat & Lift app.

## Getting Started

### Prerequisites

- Golang - https://go.dev/learn/
- Open Food Facts Data - https://world.pro.openfoodfacts.org/data

### Installation

1. Download the latest JSONL gzipped Data Export from Open Food Facts
2. Place the JSONL gzipped file in the `input` folder in the project root
3. Run the project

```console
go run main.go
```