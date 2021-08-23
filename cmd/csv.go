package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"os"
	"regexp"
)

var invalidAnswer *regexp.Regexp

func init() {
	invalidAnswer = regexp.MustCompile(`^(|.+ \/ 0)$`)
}

func isValidAnswer(answer string) bool {
	return !invalidAnswer.MatchString(answer)
}

func loadAnswers(testFile string) *[]Response {
	// Read file from testFiles
	content, err := os.ReadFile(testFile)
	if err != nil {
		log.Fatal(err)
	}

	// Check and Skip first row
	r := csv.NewReader(bytes.NewReader(content))
	_, err = r.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare result slice
	var responses = []Response{}

	for {
		// Read row
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Row processing
		datetime := row[0]
		sex := row[2]

		answers := []string{}
		for _, answer := range row[3:] {

			if isValidAnswer(answer) {
				answers = append(answers, answer)
			}
		}

		responses = append(responses, Response{datetime, sex, answers})
	}
	return &responses
}

func writeResults(fileName string, records [][]string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
