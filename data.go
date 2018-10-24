package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func saveData(content Content) {
	bytes, err := json.Marshal(content)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to marshal json")
		os.Exit(1)
	}

	ioutil.WriteFile(dataFile, bytes, 0644)
}

func loadData() Content {
	var content Content

	jsonData, err := ioutil.ReadFile(dataFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read data file : '%s'\n", dataFile)
		os.Exit(1)
	}
	err = json.Unmarshal(jsonData, &content)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to parse JSON data")
		os.Exit(1)
	}

	return content
}
